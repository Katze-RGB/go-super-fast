goal:
parse a large csv (1b+ rows) quickly, reliabily, and making effective use of system resources, and do something with the results, averages or w/e

potential constraints:
we don't have 20gb+ ram, so probably want to stream files
make effective usage of all available system resources, so do parallelism and concurrency stuff. Good place to experiment with channels.
error handling: it's gotta be decent

strategy:
- chunk csv up into multiple parts (strategy for this? Chunk number based on number of cores available? Chunk size based on ram availability?)
- golang's standard csv reader already uses bufio looks like, probably don't have to worry much about ram unless we do something really dumb
- after creating chunks, we'll get basically a 2d array of strings.
- each chunk gets blorped in a processing channel, and we do things (like what?) to each one from there
- get a 1b row csv where the first column is an integer field. I'm not uploading a 16gb csv to github lol

questions:
- how do we make sure all our goroutines actually finish and we don't exit prematurely?
  -answer: waits
- how do we avoid overflowing channel buffers? Do we actually need buffered channels?
  -we do need buffers, and we can solve it pragmatically by setting buffer size to be numWorkers*chunkSize
- would mutex be better for this than channels? Is there realistically a difference?
  -mutexes would probably work but channels are better
- is it faster to use the csv reader or do it bytewise (probably faster to do it bytes but also not worth the readability/maintanability hit?)
  -it would be but reimplementing buffered reading sounds like a pain and not very flexible
- worker strategy: is there a good rule of thumb around number of workers per thread or is it workload dependant? I'd expect for a task with a lot of io would benefit from >1 worker per thread
