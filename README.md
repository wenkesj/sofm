# Self Organized Feature Map
SOFM Associative memory written in Golang

```sh
go get github.com/wenkesj/sofm
```

# Usage
Super easy to use, just plug and chug

```sh
sofm
  --width=< 1,INF >                   # Width of the SOFM (from 1 - Infinity)
  --height=< 1,INF >                  # Height of the SOFM (from 1 - Infinity)
  --train-data=path/to/train/data     # Path of the train data file
  --train-labels=path/to/train/labels # Path of the train data labels file
  --test-data=path/to/test/data       # Path of the test data file (optional)
  --test-labels=path/to/test/labels   # Path of the test data labels file (optional)
  --learning-rate=< 0,1 >             # Learning rate of the SOFM (from 0 - 1)
  --iterations=< 1,INF >              # Number of iterations to train on (from 0 - Infinity)
  --output=path/to/output             # Path to the output file (optional)
  --save=path/to/save.gob             # Path of the saved network file (optional)
  --load=path/to/load.gob             # Load a network file (optional)
  --normalize                         # Normalize the data to -1,+1 range (optional)
```

# [Example](https://github.com/wenkesj/sofm/tree/master/example)
**SOFM** on a set of 16 animals with 13 attributes that classify them into different classes.

To run the example **without** normalized inputs

```sh
sofm --width=10 --height=10 --train-data=data/train/data.txt --train-labels=example/train/labels.txt --test-data=example/test/data.txt --test-labels=example/test/labels.txt --learning-rate=0.2 --iterations=32000 --output=results.normal.txt --save=example/save/network.gob
```

To run the example **with** normalized inputs

```sh
sofm --width=10 --height=10 --train-data=example/train/data.txt --train-labels=example/train/labels.txt --test-data=example/test/data.txt --test-labels=example/test/labels.txt --learning-rate=0.2 --iterations=32000 --output=results.thresholds.txt --save=example/save/network.gob --normalize
```
