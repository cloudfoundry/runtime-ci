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

## Release and Stemcell Updates
#{releases_table}
HEREDOC
  end

  private

  def render_table(release_updates)
    header = <<-HEADER
| Release | Old Version | New Version |
| ------- | ----------- | ----------- |
HEADER

    table = ""
    release_updates.each do |release_name, release_update|
      table << "| #{release_name} | #{release_update.old_version} | #{release_update.new_version} |\n"
    end
    "#{header}#{table}"
  end
end
