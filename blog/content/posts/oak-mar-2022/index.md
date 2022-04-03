---
title: "State of Oak: March 2022"
date: 2022-04-02T18:26:02-06:00
draft: false
---

{{< figure src="oak_with_gophers.png" width=400 height=400 >}}

## What's new this month

- Oak V4 Alpha
- Bark: A grid-based game engine

### Oak v4 Alpha

[Oak v4 alpha 1](https://github.com/oakmound/oak/releases) is out today, featuring a complete overhaul of the event system, enabled via type parameterization and go 1.18.
While we're planning a breaking change, there's several things we already know we want to include:

#### Event Overhaul

We've started off by rewriting the `event` package. There are a lot of changes, but the biggest difference is the use of type parameters to make the binding / triggering interface type safe:

Before:

```go
func() {
    ...
    k.Bind(key.Down, key.Binding(func(id event.CID, ev key.Event) int {
        kb, _ := k.ctx.CallerMap.GetEntity(id).(*Keyboard)
        btn := ev.Code.String()[4:]
        if kb.rs[btn] == nil {
            return 0
        }
        kb.rs[btn].Set("pressed")
        kb.rs[ev.Code].Set("pressed")
        return 0
    }))
}

...
package key

func Binding(fn func(event.CID, Event) int) func(event.CID, interface{}) int {
	return func(cid event.CID, iface interface{}) int {
		ke, ok := iface.(Event)
		if !ok {
			return event.UnbindSingle
		}
		return fn(cid, ke)
    }
}
```

After:

```go
event.Bind(ctx, key.AnyDown, k, func(kb *Keyboard, ev key.Event) event.Response {
    if kb.rs[ev.Code] == nil {
        return 0
    }
    kb.rs[btn].Set("pressed")
    kb.rs[ev.Code].Set("pressed")
    return 0
}))
```

In short, all the helper functions and bad feeling "reload this thing" and "assert it is what it should be" are now internal to `event` itself. Because Go 1.18 does not support type-parameterized methods, these new features are not things you can call on an `event.Handler` internally, but are functions that accept `event.Handler`s instead.

There are many other significant changes in this first release already (and with this first release we are probably done looking at `event`):

The handler interface has changed:

Before:

```go
type Handler interface {
	WaitForEvent(name string) <-chan interface{}
	// <Handler>
	UpdateLoop(framerate int, updateCh chan struct{}) error
	FramesElapsed() int
	SetTick(framerate int) error
	Update() error
	Flush() error
	Stop() error
	Reset()
	SetRefreshRate(time.Duration)
	// <Triggerer>
	Trigger(event string, data interface{})
	TriggerBack(event string, data interface{}) chan struct{}
	TriggerCIDBack(cid CID, eventName string, data interface{}) chan struct{}
	// <Pauser>
	Pause()
	Resume()
	// <Binder>
	Bind(string, CID, Bindable)
	GlobalBind(string, Bindable)
	UnbindAll(Event)
	UnbindAllAndRebind(Event, []Bindable, CID, []string)
	UnbindBindable(UnbindOption)
}
```

After:

```go
type Handler interface {
    Reset()
	TriggerForCaller(cid CallerID, event UnsafeEventID, data interface{}) chan struct{}
	Trigger(event UnsafeEventID, data interface{}) chan struct{}
	UnsafeBind(UnsafeEventID, CallerID, UnsafeBindable) Binding
	Unbind(Binding) chan struct{}
	UnbindAllFrom(CallerID) chan struct{}
    SetCallerMap(*CallerMap)
	GetCallerMap() *CallerMap
}
```

Yes, it's drastically smaller! The following realizations enabled this:

- Many of these methods (SetTick, UpdateLoop, FramesElapsed, Stop, Pause, Resume) all exist to enable looping event handlers over event.Enter calls. All of these functions, however, were either not used by oak or could be removed by changing this feature from a method on a handler to a helper function. The new function `event.EnterLoop` handles this use case.
- GlobalBind is now a type-safe helper function, calling `UnsafeBind` with the `event.Global` constant
- `TriggerBack` is now just `Trigger`, with the distinction on use of the returned channel documented
- `UnsafeBind` returns a new `Binding` type, which can be used to enable the other Unbind variants
- `SetRefreshRate` and `Flush` were needed to run `ResolveChanges`, which deferred handling of bindings to a single looping thread. Each binding now creates a goroutine, waits on a lock, takes effect once it acquires the lock, and then closes the returned `chan struct{}`. This drastically simplifies the package's internal systems.  

And the new methods:

- `Set` and `Get` for the bus's caller map, to make it easier to obtain callers from the correct map when needed.
- `UnbindAllFrom` was identified as a cleaner variant of former unbind helpers for Callers.

There was one problem we had to solve in the details above: previously the `ResolveChanges` loop was managed by two mutexes, and was designed so that concurrent calls to `Reset` on the handler would not cause bindings to operate against the newly reset bus, for example. So before:

```go
func (eb *Bus) Bind(name string, callerID CID, fn Bindable) {
	eb.pendingMutex.Unlock()
    eb.binds = append(eb.binds, UnbindOption{
        Event: Event{
			Name:     name,
			CallerID: callerID,
		}, Fn: fn})
	eb.pendingMutex.Unlock()
    // the binding appended above (confusingly called an 'UnbindOption') will be picked up next loop, or dropped next Reset.
}
```

After:

```go
func (bus *Bus) UnsafeBind(eventID UnsafeEventID, callerID CallerID, fn UnsafeBindable) Binding {
	expectedResetCount := bus.resetCount
	bindID := BindID(atomic.AddInt64(bus.nextBindID, 1))
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		bus.mutex.Lock()
		defer bus.mutex.Unlock()
		if bus.resetCount != expectedResetCount {
			// The event bus has reset while we we were waiting to bind this
			return
		}
		bl := bus.getBindableList(eventID, callerID)
		bl[bindID] = fn
	}()
	return Binding{
		Handler:       bus,
		EventID:       eventID,
		CallerID:      callerID,
		BindID:        bindID,
		Bound:         ch,
		busResetCount: bus.resetCount,
	}
}
```

Internally event's built in handler (`Bus`) tracks the number of times it has reset, and we use this within binding operations to recognize when a call has been made invalid due to a concurrent reset on the bus.

Some ancillary packages changed along with this overhaul:

- Scenes no longer have a `Loop` function, because we never used them and if you wanted one you could easily get one via `event.GlobalBind`
- `scene.Context` now embeds its event types, to make it easy to perform event operations using it.
- The `key` and `mouse` packages have a different event syntax to better differentiate between listening for e.g. any key was pressed vs the 'w' key was pressed.

There's a lot more detail in the package itself documenting why certain things are the way they are, which was also a lacking quality from the v3 event package.

#### Drivers / Shiny

In progress but not in the alpha release, we are attempting to overhaul the internal os driver interface. Goals are:

- To remove and consolidate concepts and code that we've inherited from exp/shiny which we no longer need.
- To move the current system of runtime checks for OS level features (which was originally introduced for backwards compatibility) to compile-time checks. E.g. right now if you were to call `oak.SetTopMost` when compiling to javascript it would let you, and that operation would just return an error whenever you called it. This change would demand that you move code using those sorts of operations to build tag guarded files for specific OSes if building a multi-platform game that wanted to do specific os level operations not globally supported.

The latter goal has a working windows implementation up here: https://github.com/oakmound/oak/pull/198/files, but there's more to do to make every platform follow this, and to hopefully enhance all of these drivers to support more of the features other OSes currently support.

#### Entities

To put it briefly: the `entities` package hasn't changed in a long time, and we've used it enough to throw it away and make something better with a clearer API. Its constructors take too many arguments, and it defines too many useless types with the goal of trying to express everything anyone could possibly want with the smallest data structures; this is pre-optimization, and we should instead just offer one thing that does everything you could want via an interface, like `btn` does. If a user is sad that that takes up 100 bytes instead of 40, one can always copy and delete things one does not want.

#### Audio Cleanup

As a part of the introduction of streaming audio, we needed to add it in such a way that it kept around the previous, non-streaming audio API. These APIs need to be combined. Unfortunately we cannot just only support streaming audio, because that's infeasible in JS-- you need to write and import JS modules to stream audio instead of just loading it all at once. This means we probably need to still support both streaming and non-streaming audio, cleanly indicating non-streaming audio is the only option in JS, and probably not enabling it on platforms that can support streaming audio.

This project is still being mulled over, obviously.

### Bark: A grid-based game engine

{{< figure src="bark_acorn.png" width=100 height=100 >}}

We've decided on a name for the super game engine we're building on top of Oak:`Bark`. and we're working on it. We're still in the API design stage, because we got a little swept up in Oak v4 work, but there should be more news shortly.

## What does it look like is coming up maybe before the next month has ended

We're going to finish the things mentioned above, being oak v4 and bark, and hopefully put out some demo / example games with the new APIs.

## Thanks for Reading

If you are using or thinking about using Oak do not hesitate to reach out with questions or suggestions.

## Also Who am I (AKA About the Author)

My name is Patrick Stephen in [The Physical World](https://en.wikipedia.org/wiki/Earth) and [200sc](https://github.com/200sc/) in [The Digital World](https://digimon.fandom.com/wiki/Digital_World).

I've been working on Oak since 2016, however much work at any time depending on whether I was receiving payment to work on something else instead for 40 hours. As of Apr 1 2022 the entity doing that is [strongDM](https://www.strongdm.com/).

If you'd like to reach out to me you may do so via: patrick.d.stephen@gmail.com
