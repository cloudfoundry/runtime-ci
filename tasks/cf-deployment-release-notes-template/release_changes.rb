require 'yaml'
require 'hashdiff'

class ReleaseUpdate
  attr_accessor :old_version, :new_version
end

class ReleaseUpdates
  class << self
    def load_from_files(filename, opsfile: false)
      release_candidate = load_yaml_file('cf-deployment-release-candidate', filename)
      release_candidate_releases_list = collect_releases_and_stemcells(release_candidate, opsfile: opsfile)

      master = load_yaml_file('cf-deployment-master', filename)
      master_releases_list = collect_releases_and_stemcells(master, opsfile: opsfile)

      release_updates = ReleaseUpdates.new
      changeSet = HashDiff.diff(master_releases_list, release_candidate_releases_list)
      changeSet.each do |change|
        release_updates.load_change(change)
      end

      release_updates
    end

    private

    def load_yaml_file(input_name, filename)
      filepath = File.join(input_name, filename)
      if File.exists? filepath
        file_text = File.read(filepath)
        parsed_yaml = YAML.load(file_text)

        return nil unless parsed_yaml
        parsed_yaml
      end
    end

    def collect_releases_and_stemcells(manifest, opsfile: false)
      return [] if manifest.nil? || manifest.empty?
      return manifest['releases'] + manifest['stemcells'] unless opsfile

      manifest.select do |op|
        op['type'] == 'replace' && (op['path'].start_with?('/releases') || op['path'].start_with?('/stemcells'))
      end.collect do |op|
        {
          "name" => op["value"]["name"] || op['value']['os'],
          "version" => op["value"]["version"]
        }
      end
    end
  end

  def initialize
    @updates = {}
  end

  def load_change(change)
    op = change[0]
    if op == "~"
      return
    end

    name = change[2]['name'] || change[2]['os']
    version = change[2]['version']

    release_update = @updates[name] || ReleaseUpdate.new

    if op == '+'
      release_update.new_version = version
    elsif op == '-'
      release_update.old_version = version
    end

    @updates[name] = release_update
  end

  def get_update_by_name(release_name)
    @updates[release_name]
  end

  def each
    @updates.each do |release_name, release_update|
      yield release_name, release_update
    end
  end

  def merge!(updates2)
    updates2.each do |release_name, release_update|
      @updates[release_name] = release_update
    end
  end
end
