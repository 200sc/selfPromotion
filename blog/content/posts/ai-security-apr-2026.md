---
title: "AI has another security problem"
date: 2026-04-20T20:00:00-06:00
draft: false
---

## AI has another security problem

We all know about how MCP has fundamental flaws in how it handles authorization,
how OpenClaw and its relatives open up local machines to exploits and the whole
world of LLMs is a security nightmare𐠒. And we also are now hearing that Mythos,
Anthropic's latest model, can exploit all of the security flaws in all of the
software ever and it's so impressive that it was released ahead of time to 
everyone important so they could run it themselves and protect themsevelves from
the oncoming onslaught of hackers equipped with it. 

But AI has another, more fundamental security problem. Let's say Mythos is real, 
and it costs something on the order of tens of thousands of dollars to find new 
exploits in the code of big open source projects that power millions of servers. 

How expensive do you think it is to find exploits in closed source systems? It
must be substantially less expensive: The depth of security in any system is a 
product of the number of eyes and hands that system has passed over. Each user 
of an open source system is one more possible misconfiguration that leads to a 
revelation of a problem to solve in an open library, and each year the code 
spends open is another year of eyes skimming or scanning the project identifying
improvements and flaws. 

Closed source software doesn't have any of these benefits. No one who is using 
it is going to be able to freely draw a line from a failure they notice to a 
problem in the code, and layers of support teams between the users and the 
engineers who could look at the code will further hamper open discovery of 
problems. 

No one is reading it passively other than those who are paid to, and they are 
famously not paid to keep it secure in general. Why spend time looking for 
ghosts of security problems when there's a customer who needs a new feature at 
the end of the week?

But closed source companies pay for security audits, you say? They pay companies
out the wazoo to analyze their codebases for problems passively! This is all 
true, and it's all security theater. As a security-conscious engineer myself, 
I've interally documented dozens of security flaws, hoping to be allocated to 
fix them, and seen half a dozen security teams of auditors and all the static
analysis tools in the business miss the interesting problems time and time 
again, instead finding "problems" like "you shouldn't cast from an int64 to an 
int32" or "user input from CLI flags can be provided to this API endpoint".

What does this mean for software? It means that large open source projects are
still the best place to put your trust, and it also reveals a chink in the 
armor of Schrödinger's AI bubble: if it is true that LLMs can find security 
problems in exchange for nominal funds, and if it persists that LLMs will produce
large pieces of software by themselves with minimal or no human review, then the
only programs that will be safe will be those which are not written by LLMs and
which go through extensive human review. And the best match for those 
conditions, is open source.

LLM-written code is inherently less secure by it's nature: models are trained on the 
average code i.e. insecure code, and it is produced with inherently less oversight 
than human-written code-- pattern matching to achieve a desired result; If Mythos can 
do what Anthropic says it can, then every LLM-written codebase is a sitting duck. 

---

(𐠒) see the top page of hackernews on any given day for references

Q: What if I use Mythos on my own code to make sure it isn't vulnerable?

A: Anthropic's existing code review tools are already prohibitively expensive. 
   Should Mythos be opened up, it'll be an order of magnitude more expensive 
   than Opus 4.7. There's substantial incentive for hackers to attack you, 
   because they could gain millions in ransom, but it's going to be a hard sell
   to your management to keep running tens of thousands in background scans as
   you release new updates, "just in case we find a big bug which could expose
   us to potentially higher costs"; companies building software in the age of
   Agile and LLMs do not want to spend this money on nothing (compared to what
   they could spend it on, on new features) and so they will not. 

Q: What if I use Mythos to write the code? Won't the security problems go away?

A: Assuming this is more affordable, not necessarily. Mythos' current tests are 
   presumably using prompts / a scaffolding specially targetting the identification
   of security flaws, and they were apparently run after the fact line-by-line. The 
   configuration one uses to produce code won't necessarily have the same traits,
   because LLMs do not think, they match patterns based on their inputs: if it lacks
   the same input structure which is being used in Anthropic's current marketing push
   then we lose any confidence in its security analysis capabilities. 

Q: What if Mythos is just hype?

A: Then this might be a smaller problem, but it's still a problem; LLM-code defensibility
   isn't any different, but the attacker might need more man-hours to produce the 
   exploits. 

## Who am I

If you'd like to reach out you may do so via: patrick.d.stephen@gmail.com