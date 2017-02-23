Slacker is Copyright 2015, 2016, 2017 Ed Marshall. All rights reserved. Use of
this source code is goverened by the terms of the GNU General Public License,
either version 3, or (at your option) any later version; please see the file
COPYING for more details.

Slacker is where I keep all my silly integrations for Slack.

Right now, the only thing it has is a simple slash command for looking up a
stock ticker symbol and displaying information about it back to the current
channel. It uses the Yahoo Finance API in what is probably a terrible and
non-canonical way. It can either respond inline, or asynchronously via the
recently-added `response_url` field; see the configuration file.

The only external dependency this project has right now is on BurntSushi's toml
package: https://github.com/BurntSushi/toml
