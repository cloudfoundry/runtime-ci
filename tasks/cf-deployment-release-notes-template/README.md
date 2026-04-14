# cf-deployment-release-notes-template

This task generates the cf-deployment release notes template.

## Run tests locally

### Prerequisites

- Ruby 3.2+
- Bundler installed

### From repository root

```bash
cd tasks/cf-deployment-release-notes-template
bundle install
bundle exec rspec
```

### From task directory (already there)

```bash
bundle install
bundle exec rspec
```

Notes:
- The local `.rspec` config is set to run `*_spec.rb` files in this directory.
- If dependencies changed, re-run `bundle install` before running tests.

