class OpsFileFinder
  def self.find_ops_files(repo_dir)
    ops_files_and_directories = Dir.glob(File.join(repo_dir, "operations", "*.yml"))
    opsfile_list = ops_files_and_directories.select { |file_or_directory| File.file?(file_or_directory) }

    experimental_ops_files = Dir.glob(File.join(repo_dir, "operations", "experimental", "*.yml"))
    opsfile_list += experimental_ops_files

    experimental_ops_files = Dir.glob(File.join(repo_dir, "operations", "addons", "*.yml"))
    opsfile_list += experimental_ops_files

    opsfile_list.map { |opsfile| opsfile.gsub!("#{repo_dir}/operations/", '') }
  end
end
