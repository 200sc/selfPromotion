---
title: "State of Oak: February 2022"
date: 2022-03-01T18:26:02-06:00
draft: false
---

## What is this Writing; Why am I writing it

See [Last month's blog post](https://www.200sc.dev/posts/oak-jan-2022/)

## What's new this month

- Audio System Rework
- Calculator App
- New game on Itch

### Audio System Rework

We've got an under review [PR](https://github.com/oakmound/oak/pull/191) overhauling the audio system. The primary change in this system is the switch from a audio system which requires the audio be loaded in memory to a reader/writer interface that works with in memory data or streaming data or dynamic audio data.

On Windows, we're still using DirectSound and it is still serving us very well-- it makes sense to me using these libraries how Windows became the primary target for games: they built a lot of good stable APIs that made this type of thing not that bad. On Linux we've switched to pulseaudio and it has some strange behavior, but we're planning on splitting out the system to support multiple backends before the PR is done, and alsa will likely be brought in as an alternative. Pulseaudio also works on OSX, and is much more feasible to do without diving into more C variant languages than any of Apple's built-in audio APIs.

This work should be merged in the coming days, unless a serious refactor or big new idea comes up.

### Calculator App

I wrote a [calculator app](https://github.com/200sc/oakcalc) in oak, mostly as a vehicle for working with buttons and working on a title-bar component. The titlebar component is designed to be robust and extensible but we'll have to see as we go along if that holds water. The other interesting note about this is that you may be surprised how much planning needs to go into a calculator in terms of parsing arithmetic. There are obviously shortcuts (if you're in a interpreted language, just `eval` the math given and cross your fingers the user didn't input anything crazy). [This implementation](https://github.com/200sc/oakcalc/blob/main/internal/arith/arith.go#L149) includes the usual shortcuts you might expect from an example project that didn't really want to get into language parsing but did anyway-- no floating point or big rational math, just integers, and no smart parsing of parentheses-- just pre-read entire paren bodies up front and then parse the inner elements.

A fun read, still.

### New Game on Itch

Oakmound Studio has put out a new small game that started as a space themed Oregon Trail, then realized that maybe Oregon Trail isn't that fun if you aren't in a classroom and the pun "Orion Trail" is already taken-- all the way swerving through a couple different genres until it landed at a daily picross game with sokoban controls: [Sokopic](https://oakmound.itch.io/sokopic). (Crossed fingers no one has used that name yet, as it's actually not that bad of a game title). It's not open source right now, but it's got a fun picross solver used to generate daily levels and is a pretty good time (I think (I'm gonna play it at least)).

(Oakmound Studio is myself, [ImplausiblyFun](https://github.com/Implausiblyfun), and [LightningFenrir](https://github.com/lightningfenrir))

## What does it look like is coming up maybe before the next month has ended

### Grid-based Super Engine

A lot of games, it turns out, are built around grids. Most puzzle games, a lot of rogue likes, a lot of early RPGs, tactics games, and so on. We've built two games this year, Viertris and now Sokopic, that are primarily built around grid interfaces, and I think its time we make a super engine on top of Oak that is specifically for everything you'd ever want in a grid based game.

This engine should (hopefully):

- Manage the size of the grid for you
- Manage the types of things or tiles that can go in the grid
- Manage how the things in the grid can interact with other things in the grid
- Manage how things in the grid accept user input
- Manage pausing
- Manage events

So that's something we're going to look at. The test of this is theoretically, once this super engine exists, a. does it interact cleanly with oak in a way that isn't magical; that feels like Go and like Oak and b. can we write a rougelike in it without much struggle?

### Generics

We've done some spiking with generics:

- Events are promising. We could theoretically make every event in Oak strongly typed, without breaking backwards compatibility. Relevant issue for discussion [here](https://github.com/oakmound/oak/issues/192)
- Geometry is also theoretically promising, but so far not as much as events, as much of our functions in the geometry libraries rely on specific information about the types present. Floating point epsilon equality and operations which require negative numbers are some examples.

### OSX ARM Support

This was on last month's list and it's still on this month's list-- I've done marginal work with ARM and OSX and the primary thing noticed so far is that it appears to have deprecated some metal color schemes, forcing oak to currently use BRGA where it thinks its using RGBA. Unless Apple really did drop support for RGBA encoded color this should just be some API digging, forking of metal libraries, and ~done.

### In progress bird based platformer

I started and half finished a flappy bird controlled platformer with a simple aesthetic and goal of-- just get some coins and reach the end in a low pressure environment

## Thanks for Reading

I do genuinely appreciate time spent reading this, if you are using or thinking about using Oak do not hesitate to reach out with questions or suggestions.

## Also Who am I (AKA About the Author)

(I am the same person I was last month, if you've already read about who I am)

My name is Patrick Stephen in [The Physical World](https://en.wikipedia.org/wiki/Earth) and [200sc](https://github.com/200sc/) in [The Digital World](https://digimon.fandom.com/wiki/Digital_World).

I've been working on Oak since 2016, however much work at any time depending on whether I was receiving payment to work on something else instead for 40 hours. As of Feb 1 2022 the entity doing that is [strongDM](https://www.strongdm.com/).

Earlier I've worked on smaller projects, discoverable via github, that are more scientific-- evolutionary computation, computational geometry, and extensible language theory as notable call-outs. These fall as a third priority on my list, however, behind work that keeps my family eating and work on Oak (Oak being work that roughly got me enough miniscule notoriety to provide the previous eating-providing work).

If you'd like to reach out to me you may do so via: patrick.d.stephen@gmail.com
