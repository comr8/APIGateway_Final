package consts

import "os"

var VerificationServiceURL = os.Getenv("VERIFICATION_SERVICE_URL")
var CommentServiceURL = os.Getenv("COMMENT_SERVICE_URL")
var NewsServiceURL = os.Getenv("NEWS_SERVICE_URL")
