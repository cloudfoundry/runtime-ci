module InputFileLoader
  class << self
    def load_yaml_file(input_name, filename)
      filepath = File.join(input_name, filename)
      if File.exists? filepath
        file_text = File.read(filepath)
        YAML.load(file_text)
      end
    end
  end
end
