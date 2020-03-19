class Renderer
  def render(release_updates:)
    releases_table = render_table(release_updates)
<<-HEREDOC
## Notices

## Manifest Updates

## Ops-files
### New Ops-files
### Updated Ops-files

## Other Updates

## Release Updates
_Warning: The Release Notes column only highlights noteworthy updates for each release bump. However, it is not exhaustive and we recommend you visit the actual release notes below before every upgrade._
#{releases_table}
HEREDOC
  end

  private

  def render_table(release_updates)
    header = <<-HEADER
| Release | Old Version | New Version | Release Notes |
| ------- | ----------- | ----------- | ------------- |
HEADER

    table = ""
    release_updates.each do |release_name, release_update|
      table << "| #{release_name} | #{render_version(release_update, 'old')} | #{render_version(release_update, 'new')} | |\n"
    end
    "#{header}#{table}"
  end

  def render_version(release_update, type)
    version = release_update.send(type + '_version')

    url_method = type + '_url'
    url = release_update.respond_to?(url_method) ? release_update.send(url_method) : nil

    return url ? "[#{version}](#{url})" : version
  end
end
