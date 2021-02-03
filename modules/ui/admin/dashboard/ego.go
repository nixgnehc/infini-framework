// Generated by ego.
// DO NOT EDIT

package dashboard

import (
	"fmt"
	"github.com/nixgnehc/infini-framework/modules/ui/common"
	"io"
	"net/http"
)

var _ = fmt.Sprint("") // just so that we can keep the fmt import for now
func Index(w http.ResponseWriter, r *http.Request) error {
	_, _ = io.WriteString(w, "\n\n")
	_, _ = io.WriteString(w, "\n")
	_, _ = io.WriteString(w, "\n\n")
	common.Head(w, "Dashboard", "")
	_, _ = io.WriteString(w, "\n<link rel=\"stylesheet\" href=\"/static/assets/css/tasks.css\" />\n<script src=\"/static/assets/js/jquery.sparkline.min.js\"></script>\n<script src=\"/static/assets/js/jquery.timeago.js\"></script>\n<script src=\"/static/assets/js/page/tasks.js\"></script>\n\n<META HTTP-EQUIV=\"refresh\" CONTENT=\"0;URL=/admin/console/\">\n\n")
	common.Body(w)
	_, _ = io.WriteString(w, "\n")
	common.Nav(w, r, "Dashboard")
	_, _ = io.WriteString(w, "\n\n\n\n<div class=\"tm-middle\">\n    <div class=\"uk-container uk-container-center\">\n\n        <div class=\"uk-width-1-1 uk-alert\">\n            <div class=\"uk-grid\">\n                <div class=\"uk-width-1-2\">Checker: <span id=\"checker_task_num\">N/A</span> <span class=\"dynamicsparkline\">Loading..</span></div>\n                <div class=\"uk-width-1-2\">Crawler: <span id=\"crawler_task_num\">N/A</span> <span class=\"dynamicbar\">Loading..</span></div>\n            </div>\n        </div>\n        <!--<a href=\"#my-id\" data-uk-modal>...</a>-->\n\n        <!--&lt;!&ndash; This is the modal &ndash;&gt;-->\n        <!--<div id=\"my-id\" class=\"uk-modal\">-->\n            <!--<div class=\"uk-modal-dialog uk-modal-dialog-blank\">...</div>-->\n        <!--</div>-->\n\n        <div class=\"uk-grid\" data-uk-grid-margin>\n            <div class=\"uk-width-1-1\">\n\n\n                <table id=\"tasks\" class=\"uk-table uk-table-hover uk-table-striped\" cellspacing=\"0\" width=\"100%\">\n                    <thead>\n                    <tr>\n                        <th>ID</th>\n                        <th>Path</th>\n                        <th>Title</th>\n                        <th>Size</th>\n                        <th>Updated</th>\n                    </tr>\n                    </thead>\n                    <tbody id=\"records\">\n                    </tbody>\n                </table>\n\n\n\n            </div>\n        </div>\n\n    </div></div>\n\n<script src=\"/static/assets/js/page/index.js\"></script>\n<script type=\"text/javascript\">\n     /* Sparklines can also take their values from the first argument\n     passed to the sparkline() function */\n    var myvalues = [10,8,5,7,4,4,1];\n    $('.dynamicsparkline').sparkline(myvalues);\n\n    /* The second argument gives options such as chart type */\n    $('.dynamicbar').sparkline(myvalues, {type: 'bar', barColor: 'green'} );\n\n\n    pointData=[1,2,3,4,5,4,3,2,1,1,2,3,4,5,4,3,2,1,1,2,3,4,5,4,3,2,1,1,2,3,4,5,4,3,2,1];\n    $(function() {\n\n        var sparklineLogin = function() {\n            $('.sparklines').sparkline(\n                    [ pointData ],\n                    {\n                        type: 'line',\n                        width: '100%',\n                        height: '25'\n                    }\n            );\n        };\n        var sparkResize;\n\n        $(window).resize(function(e) {\n            clearTimeout(sparkResize);\n            sparkResize = setTimeout(sparklineLogin, 500);\n        });\n        sparklineLogin();\n    });\n</script>\n\n")
	common.Footer(w)
	_, _ = io.WriteString(w, "\n")
	return nil
}
