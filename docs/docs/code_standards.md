# Code standards

## Project values

For deciding the rules laid out below, and possibly when should we break those rules, we have these values to follow for
our pet projects:

1. **Fun:** My primary goal is to have fun. Everything else is superseded. I should always ask myself if I'm having fun
   and what's stopping me from having fun.
2. **Delight:** I want to make my life easier as a developer or an artist. I want to skip to the fun part. I should use
   the tools that will help me or make those tools whenever I can.
3. **Empathy:** I pretend another developer will eventually use or look at what I’m making. This sounds weird but I find
   myself having less empathy on myself, but I highly value people working with me. If I'm struggling to value myself,
   then maybe if I pretend to treat myself to be as another person, then maybe I can show myself more empathy too.
4. **Transparency:** Share and document the successes and failures. Oh, the mortifying ordeal of being known...
5. **Focus:** This one was not even a value. It was originally "avoid premature micro-optimization" because I overthink
   too much. I need to focus on what makes me happy, and that's expressing myself through games and whatever project I
   want to make.
6. **Forgiveness:** Self-improvement is okay, but don’t be too critical and know you did your best at the time. It's
   something I've been learning these last few months, and maybe this year and the next year is meant to forgive myself
   for my shortcomings. Accept who I am. Be thankful of all the moments I tried my best, whether I failed or succeeded.

### Context

I've been reflecting about my pet projects that have either never been started or have never been finished. I'm
currently reviewing my self. Things have not been going well. And on a related note, what made most of these projects
just kinda not go well. These projects were meant to be fun, but they now scare me...

So, I wrote down what I think went wrong, and what I could do better in the future.

**What did not go well?**

1. I overthink about the solutions.
2. They're supposed to be hobbies but I'm always finding a way to turn them into the constant self-improvement that just
   disheartens me at the end of the day.
3. I don't reward myself and would often punish myself for not being on par with what I want.

I want to change the framework of my mindset. Not that I have a specific goal, but I was thinking about what values
should I uphold, so I could get back to having fun with my hobbies.

I've kinda set myself some values that I want to remind myself. Maybe I'll formulate something more specific out of
this.

## Logging rules

We came up with these rules to help debugging but also avoid clutter and info leaking to the public.

1. Always use `WrapErrorWithDetails` at the deepest level of your code. If you got an error from a code you called
   yourself, reconsider wrapping the error with this.
    - If you want to wrap an error, use our custom `ErrorWrapper`.
2. It's okay to log on every level of the stack using `LogDetailedError`.
3. Do not show the client the errors as it is! Use `ErrorWrapper.UserMessage`!
