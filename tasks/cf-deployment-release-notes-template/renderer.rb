class Renderer
  def render(release_updates:)
    header = <<-HEADER
| Release | New Version | Old Version |
| ------- | ----------- | ----------- |
HEADER

    table = ""
    release_updates.each do |release_name, release_update|
      table << "| #{release_name} | #{release_update.new_version} | #{release_update.old_version} |\n"
    end
    "#{header}#{table}"
  end
end
