### Usage

* example

  ```
  $ ./run_cmd -cmd "ls -lh" -loop 20 -delay 1s
  ```

* params

  * `-cmd` : the command to run
  * `-loop` : run times; `-1` for no limit, which means it will run forever until you kill the program
  * `-delay` : time interval,  such as:  `1s` ,   `10ms` 