package main

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	stageWrapper := func(done In, stage Stage, input In) Out {
		stageOut := stage(input)
		wrapperOut := make(Bi)

		go func() {
			defer close(wrapperOut)

			for {
				select {
				case <-done:
					return
				case value, opened := <-stageOut:
					if !opened {
						return
					}

					wrapperOut <- value
				}
			}
		}()

		return wrapperOut
	}

	for _, stage := range stages {
		in = stageWrapper(done, stage, in)
	}

	return in
}
