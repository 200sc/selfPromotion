---
title: "We need boring tech (again)"
date: 2026-11-19T07:00:00-06:00
draft: false
---

## We need boring tech (again)

The news of the moment is that the world is on fire with code exploits. From Microsoft to flights to secret government networks, everything is under attack (or falling apart from within). It's probably correct to put a lot of blame here on AI, but for fun let's imagine some things that we can actually do about this, because until the wealthy cash out on Anthropic, those who get to decide if we have jobs and can afford to eat are not going to drop or decrease AI use. 

One thing I've noticed working in a Go / Typescript monorepo for the last several years is, for every Go vulnerability that we get notified about, we get ~40 JS vulnerabilities. Literally right now I'm looking at 34 JS vulns and 5 Go vulns sitting unpatched on one repo. There's a clear reason why there's that difference in vulnerabilities: Javascript's ecosystem is much friendlier to sprawling dependency trees. 

If we followed a boring axiom like, [a little copying is better than a little dependency](https://www.youtube.com/watch?v=PAAkCSZUG1c&t=568s), we'd see a lot fewer of these. Most of these vulnerabilities are in edge case functionality which we have no intention of using; stop importing the dependency, just take what you need, and you'll be a lot happier. This is boring, mundane work. 

Another way you might see these recent security problems, is if a dependency gets taken over by a malicious actor and updated to include malware. This is a python and javascript problem which you should never see in Go, because you should be vendoring your dependencies and only updating them as needed. Again, it's not very exciting to copy your dependencies into your source code, but it makes builds more reliable and removes various kinds of security problems. 

Imagine a test for our technical decisions: what is the most boring thing we could do? 

- What database should we use? What's the most boring thing we could do? Postgres.

- How should we host our server? What's the most boring thing we could do? Docker in a stable cloud (AWS, GCP)

- How should our app be configured? Most Boring? Env variables.

- How should our server handle restarts for long running jobs? Most Boring? Just turn the server off and on and restart the whole job. When this stops working, do a little bit more.

- How big should our PRs be? Most Boring? As small as possible.

- How much should we write to document the code? Most Boring? As little as possible, the code should be understandable already.

- How big should our issues be? Only so long as they need to be to define the problem or requirements, and no longer. 

Boring, or as its sometimes referenced, simple tech isn't worse or less useful or less intelligent than other tech. [Simplicity is complicated](https://go.dev/talks/2015/simplicity-is-complicated.slide#1). [Perfection is attained not when there is nothing more to add, but when there is nothing more to remove](https://en.wikiquote.org/wiki/Antoine_de_Saint-Exup%C3%A9ry).  

Whether we're building a library, website, or server, the simpler and smaller an output we build, the easier it will be to maintain, and the more secure it will be. And now that it's easy to exchange funds for verbose AI-generated programs, simplicity and refinement too might be the only way to stand out. 

Maybe we need more boring, mundane work for a better future. 

## Who am I

If you'd like to reach out you may do so via: patrick.d.stephen@gmail.com