class VariableUpdate
  attr_accessor :state, :type
end

class VariableUpdates
  class << self
    def load_from_files(filename)
      master = InputFileLoader.load_yaml_file('cf-deployment-master', filename)
      release_candidate = InputFileLoader.load_yaml_file('cf-deployment-release-candidate', filename)

      variables_master = master['variables']
      variables_release_candidate = release_candidate['variables']

      updates = VariableUpdates.new
      changeSet = HashDiff.diff(variables_master, variables_release_candidate)
      changeSet.each do |change|
        updates.load_change(change)
      end
      updates
    end
  end

  def initialize
    @updates = {}
  end

  def load_change(change)
    op = change[0]
    name = change[2]['name']
    type = change[2]['type']

    if @updates[name]
      update = @updates[name]
      raise 'Disallowed No-op' if update.state == :added && op == '+'
      raise 'Disallowed No-op' if update.state == :removed && op == '-'

      update.state = :updated
      update.type = type
      @updates[name] = update
    else
      @updates[name] = VariableUpdate.new.tap do |o|
        o.type = type
        if op == '-'
          o.state = :removed
        else
          o.state = :added
        end
      end
    end

  end

  def get_update_by_name(variable_name)
    @updates[variable_name]
  end

  def each
    @updates.each do |variable_name, variable_update|
      yield variable_name, variable_update
    end
  end
end
