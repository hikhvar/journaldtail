# journaldtail

This is a small PoC to explore how to ship journald logs to grafana loki. This code is currently work in progress.

Any comments and suggestions are welcome. 

I made this project since the promtail community is not sure if promtail will support journald. (see:  https://github.com/grafana/loki/pull/26#issuecomment-446961639 )

## Building

To build this, you'll need cgo and the systemd headers installed - for example, on a Debian based system, `apt-get install libsystemd-dev`.

## Journal fields

When you import Journald logs to loki, you'll see a large number of log labels.
Its worth reading the [systemd journal fields documentation](https://www.freedesktop.org/software/systemd/man/systemd.journal-fields.html)
to get a feel for what is available to you.

## ToDo

- [ ] Fix logging infrastructure in code
- [ ] Enable configuration via flagext
- [ ] Support relabeling config like in promtail
- [ ] Tests
- [ ] Build and release pipeline
- [ ] Documentation  
- [ ] Store journald cursor position on disk to allow restart of journaldtail
