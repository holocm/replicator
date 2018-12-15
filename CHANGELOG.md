# v0.4.0 (2018-12-15)

Changes:

- Replace `REPLICATOR_VARIABLES_DIR` env variable by more versatile `REPLICATOR_INPUTS` variable.
- Add template function `toToml`. Works just like `toJson` from Sprig.

# v0.3.0 (2017-09-30)

Bugfixes:

- Use "html/template" instead of "text/template" to avoid unwanted quoting.

# v0.2.0 (2017-07-30)

Changes:

- Merge configuration files in a more versatile and logical way.

# v0.1.0 (2017-07-23)

Initial release. Should be functional, but very bare-bones (e.g. no
documentation except for README). Sprig version is 2.12.0.
