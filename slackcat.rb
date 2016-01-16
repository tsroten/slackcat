class Slackcat < Formula
  desc "Simple command-line Utility to post messages to Slack."
  homepage "https://github.com/vektorlab/slackcat"
  url "https://github.com/vektorlab/slackcat/archive/v0.6.tar.gz"
  version "0.6"
  sha256 "58beac16e8949a793400025ea3ce159220f21cbf3f92bf8e5530d7662d3132e9"

  depends_on "go"

  def install
    platform = `uname`.downcase.strip

    unless ENV["GOPATH"]
      ENV["GOPATH"] = "/tmp"
    end

    system "make"
    bin.install "build/slackcat-0.6-#{platform}-amd64" => "slackcat"

    puts "Ready to go! Generate a new Slack key with 'slackcat --configure'"
  end

  test do
    assert_equal(0, "/usr/local/bin/slackcat")
  end
end
