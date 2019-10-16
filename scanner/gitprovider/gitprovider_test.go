/*
 * Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
 * Use of this source code is governed by an MIT-style license that can be found in the LICENSE file
 */

package gitprovider

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var server *httptest.Server

func TestMain(m *testing.M) {
	server = setupServer()
	code := m.Run()
	defer os.Exit(code)
	teardownServer(server)
}

func setupServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		path := strings.Trim(req.URL.Path, "/")
		pathParts := strings.Split(path, "/")

		if len(pathParts) < 1 {
			_, _ = rw.Write([]byte(``))
			return
		}

		switch pathParts[0] {
		case "github":
			// https://api.github.com/repos/jquery/jquery
			_, _ = rw.Write([]byte(`{"id":167174,"node_id":"MDEwOlJlcG9zaXRvcnkxNjcxNzQ=","name":"jquery","full_name":"jquery/jquery","private":false,"owner":{"login":"jquery","id":70142,"node_id":"MDEyOk9yZ2FuaXphdGlvbjcwMTQy","avatar_url":"https://avatars1.githubusercontent.com/u/70142?v=4","gravatar_id":"","url":"https://api.github.com/users/jquery","html_url":"https://github.com/jquery","followers_url":"https://api.github.com/users/jquery/followers","following_url":"https://api.github.com/users/jquery/following{/other_user}","gists_url":"https://api.github.com/users/jquery/gists{/gist_id}","starred_url":"https://api.github.com/users/jquery/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/jquery/subscriptions","organizations_url":"https://api.github.com/users/jquery/orgs","repos_url":"https://api.github.com/users/jquery/repos","events_url":"https://api.github.com/users/jquery/events{/privacy}","received_events_url":"https://api.github.com/users/jquery/received_events","type":"Organization","site_admin":false},"html_url":"https://github.com/jquery/jquery","description":"jQuery JavaScript Library","fork":false,"url":"https://api.github.com/repos/jquery/jquery","forks_url":"https://api.github.com/repos/jquery/jquery/forks","keys_url":"https://api.github.com/repos/jquery/jquery/keys{/key_id}","collaborators_url":"https://api.github.com/repos/jquery/jquery/collaborators{/collaborator}","teams_url":"https://api.github.com/repos/jquery/jquery/teams","hooks_url":"https://api.github.com/repos/jquery/jquery/hooks","issue_events_url":"https://api.github.com/repos/jquery/jquery/issues/events{/number}","events_url":"https://api.github.com/repos/jquery/jquery/events","assignees_url":"https://api.github.com/repos/jquery/jquery/assignees{/user}","branches_url":"https://api.github.com/repos/jquery/jquery/branches{/branch}","tags_url":"https://api.github.com/repos/jquery/jquery/tags","blobs_url":"https://api.github.com/repos/jquery/jquery/git/blobs{/sha}","git_tags_url":"https://api.github.com/repos/jquery/jquery/git/tags{/sha}","git_refs_url":"https://api.github.com/repos/jquery/jquery/git/refs{/sha}","trees_url":"https://api.github.com/repos/jquery/jquery/git/trees{/sha}","statuses_url":"https://api.github.com/repos/jquery/jquery/statuses/{sha}","languages_url":"https://api.github.com/repos/jquery/jquery/languages","stargazers_url":"https://api.github.com/repos/jquery/jquery/stargazers","contributors_url":"https://api.github.com/repos/jquery/jquery/contributors","subscribers_url":"https://api.github.com/repos/jquery/jquery/subscribers","subscription_url":"https://api.github.com/repos/jquery/jquery/subscription","commits_url":"https://api.github.com/repos/jquery/jquery/commits{/sha}","git_commits_url":"https://api.github.com/repos/jquery/jquery/git/commits{/sha}","comments_url":"https://api.github.com/repos/jquery/jquery/comments{/number}","issue_comment_url":"https://api.github.com/repos/jquery/jquery/issues/comments{/number}","contents_url":"https://api.github.com/repos/jquery/jquery/contents/{+path}","compare_url":"https://api.github.com/repos/jquery/jquery/compare/{base}...{head}","merges_url":"https://api.github.com/repos/jquery/jquery/merges","archive_url":"https://api.github.com/repos/jquery/jquery/{archive_format}{/ref}","downloads_url":"https://api.github.com/repos/jquery/jquery/downloads","issues_url":"https://api.github.com/repos/jquery/jquery/issues{/number}","pulls_url":"https://api.github.com/repos/jquery/jquery/pulls{/number}","milestones_url":"https://api.github.com/repos/jquery/jquery/milestones{/number}","notifications_url":"https://api.github.com/repos/jquery/jquery/notifications{?since,all,participating}","labels_url":"https://api.github.com/repos/jquery/jquery/labels{/name}","releases_url":"https://api.github.com/repos/jquery/jquery/releases{/id}","deployments_url":"https://api.github.com/repos/jquery/jquery/deployments","created_at":"2009-04-03T15:20:14Z","updated_at":"2019-10-08T03:15:04Z","pushed_at":"2019-10-07T17:31:52Z","git_url":"git://github.com/jquery/jquery.git","ssh_url":"git@github.com:jquery/jquery.git","clone_url":"https://github.com/jquery/jquery.git","svn_url":"https://github.com/jquery/jquery","homepage":"https://jquery.com/","size":29758,"stargazers_count":52284,"watchers_count":52284,"language":"JavaScript","has_issues":true,"has_projects":true,"has_downloads":false,"has_wiki":true,"has_pages":false,"forks_count":18631,"mirror_url":null,"archived":false,"disabled":false,"open_issues_count":82,"license":{"key":"mit","name":"MIT License","spdx_id":"MIT","url":"https://api.github.com/licenses/mit","node_id":"MDc6TGljZW5zZTEz"},"forks":18631,"open_issues":82,"watchers":52284,"default_branch":"master","organization":{"login":"jquery","id":70142,"node_id":"MDEyOk9yZ2FuaXphdGlvbjcwMTQy","avatar_url":"https://avatars1.githubusercontent.com/u/70142?v=4","gravatar_id":"","url":"https://api.github.com/users/jquery","html_url":"https://github.com/jquery","followers_url":"https://api.github.com/users/jquery/followers","following_url":"https://api.github.com/users/jquery/following{/other_user}","gists_url":"https://api.github.com/users/jquery/gists{/gist_id}","starred_url":"https://api.github.com/users/jquery/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/jquery/subscriptions","organizations_url":"https://api.github.com/users/jquery/orgs","repos_url":"https://api.github.com/users/jquery/repos","events_url":"https://api.github.com/users/jquery/events{/privacy}","received_events_url":"https://api.github.com/users/jquery/received_events","type":"Organization","site_admin":false},"network_count":18631,"subscribers_count":3450}`))
		case "bitbucket":
			// https://api.bitbucket.org/2.0/repositories/litmis/mama
			_, _ = rw.Write([]byte(`{"scm":"git","website":"","has_wiki":false,"uuid":"{66d020c4-16c3-4b96-a051-a0094a100750}","links":{"watchers":{"href":"https://api.bitbucket.org/2.0/repositories/litmis/mama/watchers"},"branches":{"href":"https://api.bitbucket.org/2.0/repositories/litmis/mama/refs/branches"},"tags":{"href":"https://api.bitbucket.org/2.0/repositories/litmis/mama/refs/tags"},"commits":{"href":"https://api.bitbucket.org/2.0/repositories/litmis/mama/commits"},"clone":[{"href":"https://bitbucket.org/litmis/mama.git","name":"https"},{"href":"git@bitbucket.org:litmis/mama.git","name":"ssh"}],"self":{"href":"https://api.bitbucket.org/2.0/repositories/litmis/mama"},"source":{"href":"https://api.bitbucket.org/2.0/repositories/litmis/mama/src"},"html":{"href":"https://bitbucket.org/litmis/mama"},"avatar":{"href":"https://bytebucket.org/ravatar/%7B66d020c4-16c3-4b96-a051-a0094a100750%7D?ts=c"},"hooks":{"href":"https://api.bitbucket.org/2.0/repositories/litmis/mama/hooks"},"forks":{"href":"https://api.bitbucket.org/2.0/repositories/litmis/mama/forks"},"downloads":{"href":"https://api.bitbucket.org/2.0/repositories/litmis/mama/downloads"},"issues":{"href":"https://api.bitbucket.org/2.0/repositories/litmis/mama/issues"},"pullrequests":{"href":"https://api.bitbucket.org/2.0/repositories/litmis/mama/pullrequests"}},"fork_policy":"allow_forks","name":"mama","project":{"key":"IIOSP","type":"project","uuid":"{ddb5632c-37eb-4e68-b5c4-bcd0d3953234}","links":{"self":{"href":"https://api.bitbucket.org/2.0/teams/litmis/projects/IIOSP"},"html":{"href":"https://bitbucket.org/account/user/litmis/projects/IIOSP"},"avatar":{"href":"https://bitbucket.org/account/user/litmis/projects/IIOSP/avatar/32"}},"name":"Open Source Projects for IBM i"},"language":"c","created_on":"2017-08-01T16:47:57.614836+00:00","mainbranch":{"type":"branch","name":"master"},"full_name":"litmis/mama","has_issues":true,"owner":{"username":"litmis","display_name":"litmis","type":"team","uuid":"{f6c9fd02-930e-489e-993c-d96793cd67f6}","links":{"self":{"href":"https://api.bitbucket.org/2.0/teams/%7Bf6c9fd02-930e-489e-993c-d96793cd67f6%7D"},"html":{"href":"https://bitbucket.org/%7Bf6c9fd02-930e-489e-993c-d96793cd67f6%7D/"},"avatar":{"href":"https://bitbucket.org/account/litmis/avatar/"}}},"updated_on":"2017-08-18T13:58:56.243437+00:00","size":1002631,"type":"repository","slug":"mama","is_private":false,"description":""}`))
		case "gitlab":
			// https://gitlab.com/api/v4/projects/7824084
			_, _ = rw.Write([]byte(`{"id":7824084,"description":"Augur - Prediction Market Protocol and Client","name":"augur","name_with_namespace":"augurproject / augur","path":"augur","path_with_namespace":"augurproject/augur","created_at":"2018-08-08T22:32:40.106Z","default_branch":"master","tag_list":[],"ssh_url_to_repo":"git@gitlab.com:augurproject/augur.git","http_url_to_repo":"https://gitlab.com/augurproject/augur.git","web_url":"https://gitlab.com/augurproject/augur","readme_url":"https://gitlab.com/augurproject/augur/blob/master/README.md","avatar_url":null,"star_count":0,"forks_count":1,"last_activity_at":"2019-10-08T06:59:56.524Z","namespace":{"id":3039570,"name":"augurproject","path":"augurproject","kind":"group","full_path":"augurproject","parent_id":null,"avatar_url":"/uploads/-/system/group/avatar/3039570/Augur-Mark-Icon-400x400.png","web_url":"https://gitlab.com/groups/augurproject"}}`))
		default:
			_, _ = rw.Write([]byte(``))
		}
	}))
}

func teardownServer(s *httptest.Server) {
	s.Close()
}
