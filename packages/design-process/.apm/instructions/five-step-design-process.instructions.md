---
description: Elon Musk's five-step engineering design process — make requirements less dumb, delete the part, optimize, accelerate, automate. Must be executed in exact order.
applyTo: "**"
---

# Elon Musk's Five-Step Engineering and Design Process

Elon Musk has frequently detailed a specific, five-step engineering and manufacturing philosophy that drives the development cycles at both SpaceX and Tesla. Perhaps most famously articulated during an extensive interview at the Starbase facility with *Everyday Astronaut* (Tim Dodd), this methodology is designed to prevent engineering bloat, eliminate inefficiencies, and drive rapid innovation. 

Musk insists that these steps **must be executed in exact order**. A common trap in engineering is to optimize or automate a process that fundamentally should not exist in the first place.

Below is a deep dive into the five steps, the reasoning behind them, and real-world examples of how they are applied.

---

## Step 1: Make Your Requirements Less Dumb

**The Philosophy:**
Every requirement or constraint in an engineering project is flawed to some degree. Musk emphasizes that requirements should never be treated as sacred, especially if they come from a "smart person." In fact, requirements from smart people are the most dangerous because teams are less likely to question them. Every requirement must come with the name of the specific person who created it, not just a department (e.g., "The Legal Department" or "Safety Team"), so that the requirement can be interrogated, challenged, and updated.

**The Execution:**
* Identify the exact human being responsible for the requirement.
* Question the underlying assumptions. Is the safety margin too high? Is it based on outdated physics or legacy industry standards?
* Continuously revise the requirements to make them "less dumb."

**Real-World Example:**
At SpaceX, engineers might face a requirement that a certain component must withstand a specific vibrational frequency. Instead of blindingly engineering the part to meet that standard, they must ask the person who set the standard *why* it was set there. Often, they discover the standard was a legacy carry-over from an entirely different rocket program and is irrelevant to the current vehicle.

## Step 2: Try Very Hard to Delete the Part or Process

**The Philosophy:**
If you are not being forced to add parts back into the design at least 10% of the time, you are not deleting enough parts. Engineers have a natural bias toward adding components "just in case." This leads to bloated, heavy, and complex systems. The most reliable, cheapest, and lightest part is the one that doesn't exist.

**The Execution:**
* Aggressively challenge the necessity of every single component, sensor, fastener, and manufacturing step.
* Accept that you will sometimes make a mistake and delete something necessary. The cost of occasionally having to add a part back is lower than the cost of carrying useless parts in every vehicle forever.

**Real-World Example:**
During the Tesla Model 3 "production hell," engineers were struggling with a robot that was failing to place a fiberglass mat over the battery packs. When Musk asked what the mat was for, the battery safety team said it was for noise reduction, while the noise reduction team said it was for battery fire safety. It turned out *neither* team actually needed it. The mat was entirely redundant and was deleted from the design, immediately solving the manufacturing bottleneck.

## Step 3: Simplify or Optimize

**The Philosophy:**
Only after you have ensured the requirements are logical (Step 1) and deleted every unnecessary part (Step 2), should you attempt to optimize the design. The most common error in engineering is spending vast amounts of time optimizing a part or process that should have been deleted in Step 2. 

**The Execution:**
* Streamline the remaining parts. 
* Consolidate multiple parts into a single casting or component.
* Reduce weight, lower costs, and simplify the manufacturing operations.

**Real-World Example:**
Tesla's introduction of "Giga Presses" allowed them to cast the entire front or rear underbody of a car as a single piece of aluminum. Previously, this section of the car consisted of 70+ individually stamped metal parts that had to be welded and riveted together. By simplifying the structural design into one piece, they optimized the manufacturing line heavily.

## Step 4: Accelerate Cycle Time

**The Philosophy:**
Once you have the right requirements, the minimum number of parts, and an optimized design, you should figure out how to do it faster. Time is the ultimate currency in innovation. Moving faster allows for more iterations, more data collection, and quicker discovery of flaws.

**The Execution:**
* Find bottlenecks in the engineering, testing, or manufacturing pipeline and widen them.
* Reduce the time between designing a part and physically testing it.
* **Warning:** Do not accelerate cycle time before completing Steps 1-3. Digging your own grave faster is not a productive strategy.

**Real-World Example:**
SpaceX's Starbase development model is built on rapid iteration. Instead of spending years building one perfect prototype, they build multiple prototypes simultaneously. If a Starship explodes on the pad or during flight, the next iteration is already rolling out of the assembly building. This accelerated cycle time means they learn from physical reality much faster than legacy aerospace companies.

## Step 5: Automate

**The Philosophy:**
Automation is the final step, not the first. You should only introduce robots and automated code after you have stripped the process down to its absolute bare minimum and optimized it for speed. 

**The Execution:**
* Introduce robotic arms, automated software pipelines, and machine learning only to processes that have survived Steps 1 through 4.
* Remove human labor from repetitive, finalized tasks to scale up production.

**Real-World Example (A Cautionary Tale):**
Musk frequently cites his own failure with the early Model 3 production line as the ultimate lesson in why Automation must come last. He attempted to build an "Alien Dreadnought" factory that was entirely automated from the start. They tried to automate the installation of fiberglass mats that shouldn't have existed (failed Step 2), and automated complex routing of wire harnesses that could have been simplified (failed Step 3). Ultimately, Tesla had to rip out millions of dollars of robotics and put humans back on the line to figure out the simplified process before re-automating.

---

## Conclusion

Elon Musk's five-step process is fundamentally an exercise in first-principles thinking and rigorous discipline. It fights human nature, which naturally gravitates toward adding complexity to solve problems. By forcing teams to aggressively subtract and simplify *before* they optimize and automate, this framework allows organizations to build radically innovative hardware at unprecedented speeds.
