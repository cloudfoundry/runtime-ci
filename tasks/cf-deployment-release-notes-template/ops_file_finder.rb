class OpsFileFinder
  def self.find_ops_files(repo_dir)
    folders_to_exclude = %r{/workaround/}
    ops_files_and_directories = Dir.glob(
      File.join(repo_dir, 'operations', '**', '*.yml')
    ).grep_v folders_to_exclude

    opsfile_list = ops_files_and_directories.select { |fd| File.file?(fd) }

    opsfile_list.map { |opsfile| opsfile.gsub!("#{repo_dir}/operations/", '') }
  end
end
