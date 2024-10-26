package reddit

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	editSuccessResponseBody = `{
  "jquery": [
    [
      0,
      1,
      "call",
      [
        "body"
      ]
    ],
    [
      1,
      2,
      "attr",
      "find"
    ],
    [
      2,
      3,
      "call",
      [
        ".status"
      ]
    ],
    [
      3,
      4,
      "attr",
      "hide"
    ],
    [
      4,
      5,
      "call",
      []
    ],
    [
      5,
      6,
      "attr",
      "html"
    ],
    [
      6,
      7,
      "call",
      [
        ""
      ]
    ],
    [
      7,
      8,
      "attr",
      "end"
    ],
    [
      8,
      9,
      "call",
      []
    ],
    [
      0,
      10,
      "call",
      [
        "body\u003Ediv.content"
      ]
    ],
    [
      10,
      11,
      "attr",
      "replace_things"
    ],
    [
      11,
      12,
      "call",
      [
        [
          {
            "kind": "t1",
            "data": {
              "subreddit_id": "t5_1pkf5",
              "approved_at_utc": null,
              "author_is_blocked": false,
              "comment_type": null,
              "edited": 1729911254.0,
              "mod_reason_by": null,
              "banned_by": null,
              "ups": 1,
              "num_reports": null,
              "author_flair_type": "text",
              "total_awards_received": 0,
              "subreddit": "StandUpComedy",
              "author_flair_template_id": null,
              "likes": true,
              "replies": "",
              "user_reports": [],
              "saved": false,
              "id": "afl392f",
              "banned_at_utc": null,
              "mod_reason_title": null,
              "gilded": 0,
              "archived": false,
              "collapsed_reason_code": null,
              "no_follow": false,
              "author": "dummy",
              "can_mod_post": false,
              "created_utc": 1729911254.0,
              "send_replies": true,
              "parent_id": "t1_jtili29",
              "score": 1,
              "author_fullname": "t2_je2ek47kj",
              "report_reasons": null,
              "approved_by": null,
              "all_awardings": [],
              "collapsed": false,
              "body": "foobarbaz",
              "awarders": [],
              "top_awarded_type": null,
              "author_flair_css_class": null,
              "author_patreon_flair": false,
              "downs": 0,
              "author_flair_richtext": [],
              "is_submitter": false,
              "body_html": "\u003Cdiv class=\"md\"\u003E\u003Cp\u003Efoobarbaz\u003C/p\u003E\n\u003C/div\u003E",
              "removal_reason": null,
              "collapsed_reason": null,
              "associated_award": null,
              "stickied": false,
              "author_premium": false,
              "can_gild": false,
              "gildings": {},
              "unrepliable_reason": null,
              "author_flair_text_color": null,
              "score_hidden": false,
              "permalink": "/r/FooBarBaz/comments/1ac7e24/dummies_for_dummy/afl392f/",
              "subreddit_type": "public",
              "locked": false,
              "name": "t1_afl392f",
              "created": 1729911254.0,
              "author_flair_text": null,
              "treatment_tags": [],
              "rte_mode": "richtext",
              "link_id": "t3_3ga9e56",
              "subreddit_name_prefixed": "r/FooBarBaz",
              "controversiality": 0,
              "author_flair_background_color": null,
              "collapsed_because_crowd_control": null,
              "mod_reports": [],
              "mod_note": null,
              "distinguished": null
            }
          }
        ],
        true,
        true,
        false
      ]
    ],
    [
      0,
      13,
      "call",
      [
        "body\u003Ediv.content .link .rank"
      ]
    ],
    [
      13,
      14,
      "attr",
      "hide"
    ],
    [
      14,
      15,
      "call",
      []
    ]
  ],
  "success": true
}`

	editRateLimitErrorBody = `{
  "jquery": [
    [
      0,
      1,
      "call",
      [
        "body"
      ]
    ],
    [
      1,
      2,
      "attr",
      "find"
    ],
    [
      2,
      3,
      "call",
      [
        ".status"
      ]
    ],
    [
      3,
      4,
      "attr",
      "hide"
    ],
    [
      4,
      5,
      "call",
      []
    ],
    [
      5,
      6,
      "attr",
      "html"
    ],
    [
      6,
      7,
      "call",
      [
        ""
      ]
    ],
    [
      7,
      8,
      "attr",
      "end"
    ],
    [
      8,
      9,
      "call",
      []
    ],
    [
      1,
      10,
      "attr",
      "find"
    ],
    [
      10,
      11,
      "call",
      [
        ".error.RATELIMIT.field-ratelimit"
      ]
    ],
    [
      11,
      12,
      "attr",
      "show"
    ],
    [
      12,
      13,
      "call",
      []
    ],
    [
      13,
      14,
      "attr",
      "text"
    ],
    [
      14,
      15,
      "call",
      [
        "Looks like you've been doing that a lot. Take a break for 3 seconds before trying again."
      ]
    ],
    [
      15,
      16,
      "attr",
      "end"
    ],
    [
      16,
      17,
      "call",
      []
    ]
  ],
  "success": false
}`
)

func TestEditResponse_Success(t *testing.T) {
	var resp EditResponse
	err := json.Unmarshal([]byte(editSuccessResponseBody), &resp)
	require.NoError(t, err)
	require.True(t, resp.Success)
}

func TestEditResponse_IsRateLimited(t *testing.T) {
	var resp EditResponse
	err := json.Unmarshal([]byte(editRateLimitErrorBody), &resp)
	require.NoError(t, err)
	require.False(t, resp.Success)
	require.True(t, resp.IsRateLimited())
}
