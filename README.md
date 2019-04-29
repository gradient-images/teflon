# Teflon

Teflon is a post-production management framework for VFX heavy cinematic shows,
like feature films, and TV series. Teflon is under early development, the code and
the documentation only published to aid discussion about the design. **Please don't
use it in production in any circumstances.**

Teflon is rather different from the pipeline frameworks we know about. Teflon is
file-system based, there is no database backend behind the logic. All metadata
is stored in the file-system together with the data they refer to. The `teflon`
tool creates a view in RAM about the file-system, makes modification on the
objects then write the results to the file-system.

There are some useful information scattered throughout the [Teflon Wiki](https://github.com/gradient-images/teflon/wiki),
but a proper white paper is in the works. Also feel free to come and chat with
us at the [Teflon Gitter community.](https://gitter.im/teflon-ppp/community)

[![GoDoc](https://godoc.org/github.com/gradient-images/teflon?status.svg)](https://godoc.org/github.com/gradient-images/teflon)


Lubricating post-production pipelines. All documentation for the project currently
lives in the [Teflon Wiki](https://github.com/gradient-images/teflon/wiki).
