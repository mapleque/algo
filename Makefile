
.PHONY: train
train:
	go run train/main.go

.PHONY: valid
valid:
	go run valid/main.go

.PHONY: play
play:
	go run play/main.go

.PHONY: dryrun
dryrun:
	go run dryrun/main.go
