require 'yaml'
require 'hashdiff'

class ReleaseUpdate
  attr_accessor :old_version, :new_version
end

class ReleaseUpdates
  class << self
    def load_from_files(filename, opsfile: false)
      release_candidate_file = File.join('cf-deployment-release-candidate', filename)
      if File.exists? release_candidate_file
        release_candidate_text = File.read(release_candidate_file)
        release_candidate = YAML.load(release_candidate_text)
      else
        release_candidate = nil
      end

      master_file = File.join('cf-deployment-master', filename)
      if File.exists? master_file
        master_text = File.read(master_file)
        master = YAML.load(master_text)
      else
        master = nil
      end

      master_releases_list = collect_releases_and_stemcells(master, opsfile: opsfile)
      release_candidate_releases_list = collect_releases_and_stemcells(release_candidate, opsfile: opsfile)

      release_updates = ReleaseUpdates.new
      changeSet = HashDiff.diff(master_releases_list, release_candidate_releases_list)
      changeSet.each do |change|
        release_updates.load_change(change)
      end

      release_updates
    end

    private

    def collect_releases_and_stemcells(manifest, opsfile: false)
      return [] unless manifest
      return manifest['releases'] + manifest['stemcells'] unless opsfile

      manifest.select do |op|
        op['type'] == 'replace' && (op['path'] == '/releases/-' || op['path'] == '/stemcells/-')
      end.collect do |op|
        {
          "name" => op["value"]["name"],
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
end
