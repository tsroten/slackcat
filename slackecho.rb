class Slackecho < Formula
  desc "Simple command-line utility to post messages to Slack."
  homepage "https://github.com/tsroten/slackecho"
  url "https://github.com/tsroten/slackecho/archive/v1.0.tar.gz"
  version "1.0"
  sha256 "e18db07eecffc180ec32d3c19378d2db0cb50dad09da695ea34d5ebafffaf53c"

  depends_on "go"

  def install
    platform = `uname`.downcase.strip

    unless ENV["GOPATH"]
      ENV["GOPATH"] = "/tmp"
    end

    system "make"
    bin.install "build/slackecho-1.0-#{platform}-amd64" => "slackecho"

    puts "Slackecho is installed. If you don't already have Slackcat,\n"
    puts "install and configure Slackcat with:"
    puts "  brew install slackcat"
  end

  test do
    assert_equal(0, "/usr/local/bin/slackecho")
  end
end
