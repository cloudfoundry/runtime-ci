class ReleaseUpdate
  attr_accessor :old_version, :new_version
end

class ReleaseUpdates
  def initialize
    @updates = {}
  end

  def load_change(change)
    op = change[0]
    name = change[2]['name']
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
