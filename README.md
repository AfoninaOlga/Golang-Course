# Golang-Course
## Task 1 — Stemming
### Build
To build the project use one of the following commands:
```bash
make
```
```bash
make build
```
```bash
go build -o myapp ./stem
```
### Run
```bash
./myapp -s <string for stemming>
```
Example
```bash
./myapp -s "i'll follow you as long as you are following me"
```
Output:
```
follow long
```
