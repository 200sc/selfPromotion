---
title: "State of Oak: January 2022"
date: 2022-02-01T18:26:02-06:00
draft: false
---

## What is this Writing; Why am I writing it

### What

This blog post is meant to be the beginning of a series of posts covering the [Oak game engine](https://www.github.com/oakmound/oak):

- What is Oak?
- Where Oak was as of the end of 2021?
- What happened to Oak during January 2022?
- What is happening with Oak going into Feburary 2022?

Followup blog posts will probably not include those first two points, but will include the latter two, with the month incrementing.

Blog post Version: 1.0.0. No edits since publishing.

### Why

This blog post serves many purposes:

#### For outsiders to Oak, this writing wants you to look at Oak and try it out

Its a game engine-- its written in Go and its different from other Go game engines because its primary mission is to **just** be written in Go-- right now you can't avoid C dependencies for targets like OSX or Android, but Windows, WASM and Linux are good and stable without it and if we can find a way to get rid of C we will pursue it. Hopefully that's interesting enough to give it a look.

#### For those few familiar with Oak, this writing wants to remind you about Oak and feel good about its progress

I don't have anything to put here, scroll down and go read the progress!

#### For the author of this writing, this writing wants to serve as motivation to continue working and advancing your craft

Good on you for writing this blog post. Write another one next month, it will be better written, edited, and contain more interesting things to share because you'll have done a lot of good work in February. Keep up the good work, but don't neglect reading books or studying languages or spending time with your family or doing the dishes.

## Oak <a name="oak"/>

### Oak in 2021 <a name="oak-in-2021"/>

At the end of 2021 Oak's most recent [milestones](https://github.com/oakmound/oak/releases) were:

1. Oak 3.0.0 in September: Multi Window
2. Oak 3.1.0 in October: WASM+JS
3. Oak 3.2.0 In December: Rendering optimizations

We'd overhauled the APIs of almost every package in Oak to remove global state, sort of. If one was to completely remove global state from the engine you'd turn code like this:

```go
txt := render.NewText("Hello World", x, y) 
render.Draw(txt)
```

Into:

```go
var ctx scene.Context
fnt, err := render.NewFont({Size:10,Color:color.RGBA{255,255,255,255},DPI:72.0})
if err != nil {
    logrus.WithErr(err).Error("failed to create basic font")
    return err 
}
txt := fnt.NewText("Hello World", x, y)
ctx.DrawStack.Draw(txt)
```

So Oak does both: it provides a `render.DefaultFont` function and exports all the methods on the font as functions of the `render` package, for times when you just want to use the basic font, and likewise moves functions from `Default$Foo` objects or constructors to the top of most packages, enabling both ease of use and specialized use for when programs start to get bigger and needs become more specific.  

The mentioned change is the primary enabler of multi-window work; where before all operations were performed on 'the window' or 'the scene that is running', now Oak programs can be structured to accept contexts which contain information one needs to operate on 'the scene that this function was called by'.

WASM and JS was a long sought after goal that other Go game engines had achieved much sooner. Oak's reluctance to use OpenGL / WebGL for fear of C dependencies made this conversion more difficult, or maybe just made it seem more difficult enough to slow work towards this milestone. Several years ago Oak had functioning JS builds, but they were amazingly slow on every web browser. The addition of optimizations, moving from GopherJS to Go's built in WASM targets, respecting web browser's animation request rates, and the introduction of Multi threading in WASM all enabled enough speed ups for Oak to run at a reasonable clip in Firefox and Chrome (and maybe other browsers as well but who's counting).

To prove that we had improved the speed of Oak and the javascript backend for a fairly real game, we built out [Dash King](https://oakmound.itch.io/dashking), roughly fully releasing it on December 22nd.

(We is myself, [ImplausiblyFun](https://github.com/Implausiblyfun), and [LightningFenrir](https://github.com/lightningfenrir))

### Oak in January 2022 <a name="oak-in-january-2022"/>

#### Grove

We compiled together useful packages from previous Oak projects and built the [oakmound/grove](https://www.github.com/oakmound/grove) package from them-- this utility and component library is meant to bootstrap new projects and we hope to continually add more to it as more common needs are found in Oak programs.

#### Game Template

Along the same lines, We (mostly ImplausiblyFun) put together [a template library](https://github.com/oakmound/game-template) to build new games from-- it comes with integrations for itch and github actions so new builds will be deployed in a playable state with built, downloadable binaries, and sets up packages for some common needs we've found. The idea is to fork (or 'use this template' as github shows it) the repo and tear it apart to meet your project's needs.

#### Viertris

With that above game template I put together [a Tetris Clone](https://www.github.com/200sc/viertris), also [on itch](https://twohundredscythes.itch.io/viertris). I'm pretty proud of how this came together and how short the code is for a pretty complete rendition of Tetris. (I assume that someone else already named their Tetris clone Viertris but have not looked it up)

#### Windows ARM Support

By [incrementing a library](https://github.com/oakmound/oak/pull/184) we added ARM support for Windows! Compared to how hard it would have been to do it if we needed to compile C and figure out cross compilers / get a GCC targetting ARM on a windows ARM machine, I think that this demonstrates some value in avoiding C dependencies.

#### Linux ARM Support

For Linux ARM we had to do a bit more work, [forking a library to add in ARM code](https://github.com/oakmound/oak/pull/182), and proving that what we did worked was another adventure-- we integrated a self hosted runner with github actions that executes every example in Oak on a linux ARM machine demonstrating that each example opens a window and terminates without crashing, and this runner now runs on every build.

We can do more with these runners-- we can add runners for osx, windows, and linux, amd and arm, fairly easily; and maybe even JS? Maybe mobile in the future?

{{< youtube M0wIq-k2swk >}}
### In Progress <a name="in-progress"/>

#### Android Support

Android support is almost fully complete, with all tried examples functioning. The branch `feature/android` contains the work in progress, but needs some more work, specifically readme writing for how to actually set up Go to build for Android in 2022 and how to test those built android APKs, etc.

{{< youtube aMEzjtXgbV4 >}}

#### Audio Overhaul

The audio system underneath oak as it lives currently works pretty OK on windows, for WAV or MP3 files, and has some difficulties on linux that [need to be investigated further](https://github.com/oakmound/oak/issues/131). Audio on OSX is non-existant, audio on JS and Android is non-existant, and the interface for using audio and using code-generated audio is rough.

We may run into barriers with breaking changes here, but we should try to get as far as we can adding and improving existing support. 

### On the Horizon <a name="on-the-horizon"/>

#### Btn Refactor

Oak currently serves an internal experimental package under `entities/x` -- this package is not beholden to Oak's semantic versioning guarantees, and there's a lot of specialized work in there that should probably be moved out to the grove now.

The most used of these packages is `btn`, a fully functional-optionized package for building user interfaces by building tiers of helpers on top of Oak. It needs to grow however, and it is probably growing too big to live right inside of Oak itself. It may even need a big enough pot that it should be `oakmound/btn` instead of `oakmound/grove/btn`. This is mostly rambling, but the short of it is the entities and sub packages beneath it need re-evaluating for where they should live.

#### Another Month-Long Game

We've done Tetris, which honestly could have been done in about a week or a day with enough high speed music and questionably healthy sources of body fuel (and / or negligence of other daily tasks). What's next? We could maybe make a Mario 1 or Legend of Zelda 1 engine with a level or two without much comittment. If we're looking to make something that produces reusable components, something grid based like Sokobon or Rogue might lend itself to producing pathing helpers, for example.

This is an open question still. Nothing comes to mind immediately as obviously the next example game to make after Tetris.

#### More Utilities

I almost finished writing a Calculator before I had to solve parsing math as a language-- I'll probably finish polishing that up and release it this month, and in addition it'd be good to get more app or utility-type programs out there. A text editor or simple web browser is probably too much of a lift; a SQL editor would require also being a text editor. A graphviz or dot viewer / editor would be fairly reasonable, but maybe of limited use. A timer app might be too simple.

Also an open question. One way or another, more small (but not miniscule) examples will be built.

#### Generics

Go 1.18 is introducing Generics / Parametric Polymorphism (roughly), and Oak should look to see how it could refactor itself to take advantage of these changes. Without digging too deep, the event system currently requires `interface{}` types and casting to known types within binding functions-- can we remove this entirely? Can `Bind` calls take in a paramaterized type that would move these checks to compile time?

Based on my understanding of generics I'm honestly not optimistic about this-- methods cannot be paramaterized and we do not want to associate a single type with the event bus / event handler construct, but maybe we can do something clever and hide that cleverness in a way that it doesn't make the library difficult to use.

#### OSX ARM Support

Apple is moving all of their Macbooks to ARM so everyone has to follow suite and quick rebuild everything to support ARM64 and AMD64-- this has to be spiked out. It may be trivial, it may be that Apple has (again) deprecated an important interface we relied on and this is a big task.

#### Oak 4?

We'll probably release Oak 4 this year, adjusting the API in a smaller way than we did in Oak 3. The main targets of this refactor would be the audio system, which was almost untouched from Oak 2 to Oak 3, and the driver interface-- for backwards compatibility new driver features are added in compile-time safe ways, but we can and should adjust those interfaces now to make them compile time guarantees that, yes, you can call `.SetFullscreen(true)` on this operating system but if you try building that for JS or Android that's not going to fly.

#### Oak 5?

If we get enough velocity we may make more major versions of Oak this year but that would require that we

1. Get to Oak 4 and
1. Recognize things in Oak 4 that need breaking changes

No API is ever perfect so whether or not we get as far as another major release all depends on how much time and how many eyes Oak gets.

## Thanks for Reading

I do genuinely appreciate time spent reading this, if you are using or thinking about using Oak do not hesitate to reach out with questions or suggestions.

## Also Who am I (AKA About the Author) <a name="also-who-am-i-aka-about-the-author"/>

My name is Patrick Stephen in [The Physical World](https://en.wikipedia.org/wiki/Earth) and [200sc](https://github.com/200sc/) in [The Digital World](https://digimon.fandom.com/wiki/Digital_World).

I've been working on Oak since 2016, however much work at any time depending on whether I was receiving payment to work on something else instead for 40 hours. As of Feb 1 2022 the name of that entity is [strongDM](https://www.strongdm.com/).

Earlier I've worked on smaller projects, discoverable via github, that are more scientific-- evolutionary computation, computational geometry, and extensible language theory as notable callouts. These fall as a third priority on my list, however, behind work that keeps my family eating and work on Oak (Oak being work that roughly got me enough miniscule notoriety to provide the previous eating-providing work).

If you'd like to reach out to me you may do so via: patrick.d.stephen@gmail.com
