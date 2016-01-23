class Slackecho < Formula
  desc "Simple command-line utility to post messages to Slack."
  homepage "https://github.com/tsroten/slackecho"
  url "https://github.com/tsroten/slackecho/archive/v1.1.tar.gz"
  version "1.1"
  sha256 "2d70f7cee38668fd94f5bf53f3333d38d63235f395795023981b220ac3dae23a"

  depends_on "go"
  depends_on "slackcat"

  def install
    platform = `uname`.downcase.strip

    unless ENV["GOPATH"]
      ENV["GOPATH"] = "/tmp"
    end

    system "make"
    bin.install "build/slackecho-1.1-#{platform}-amd64" => "slackecho"

    puts "Slackecho is installed. If you don't already have Slackcat,\n"
    puts "install and configure Slackcat with:"
    puts "  brew install slackcat"
  end

  test do
    assert_equal(0, "/usr/local/bin/slackecho")
  end
end
