from(db:"testdb")
  |> range(start: 2018-05-22T19:53:26Z)
  |> filter(fn: (r) => r._measurement == "diskio" and r._field == "io_time")
  |> group(by: ["_measurement", "name"])
  |> max()
  |> map(fn: (r) => {_time: r._time, max: r._value})
  |> yield(name: "0")
