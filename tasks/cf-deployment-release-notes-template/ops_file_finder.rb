class OpsFileFinder
  def self.find_ops_files(repo_dir)
    ops_files_and_directories = Dir.glob(
      File.join(repo_dir, 'operations', '**', '*.yml')
    )

    opsfile_list = ops_files_and_directories.select do |fd|
      File.file?(fd) unless fd.match(/(use-compiled-releases.yml|test\/fips-stemcell.yml)/)
    end

    opsfile_list.map { |opsfile| opsfile.gsub!("#{repo_dir}/operations/", '') }
  end
end
