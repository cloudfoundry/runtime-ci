class ReleaseTeam
  attr_reader :anchor_github, :anchor_slack
  attr_reader :pm_github, :pm_slack
  attr_reader :name

  def initialize(name:, anchor_github:, anchor_slack:, pm_github:, pm_slack:)
    @name = name
    @anchor_github = anchor_github
    @anchor_slack = anchor_slack
    @pm_github = pm_github
    @pm_slack = pm_slack
  end

  def slack_handles
    return @anchor_slack, @pm_slack
  end

  def github_handles
    return @anchor_github, @pm_github
  end
end



class ReleaseTeamCollection
  attr_reader :teams

  def initialize(*args)
    @teams = []
    args.each do |arg|
      raise "Not a ReleaseTeam" unless arg.is_a? ReleaseTeam
      teams << arg
    end
  end

  def find_team_by_github_handles(anchor_github, pm_github)
    teams.each do |team|
      if matches_name?(team.anchor_github, anchor_github) && matches_name?(team.pm_github, pm_github)
        return team
      end
    end
    nil
  end

  private

  def matches_name?(name1, name2)
    if name1.nil?
      name1 = ''
    end

    if name2.nil?
      name2 = ''
    end

    return name1 == name2
  end
end
