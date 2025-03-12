// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.833

package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import (
	"fmt"

	"github.com/a-h/templ"
	templruntime "github.com/a-h/templ/runtime"

	"gitlab.com/hmajid2301/gofeedback/internal/transport/http/auth"
	"gitlab.com/hmajid2301/gofeedback/internal/transport/http/views/icons"
)

func NavBar(auth auth.AuthState) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<div><navbar class=\"py-2 px-4 shadow-sm navbar bg-base-100\"><div class=\"flex-1\"><a href=\"/\" class=\"inline-flex h-auto hover:bg-opacity-20 btn btn-ghost hover:btn-ghost\"><img class=\"w-auto h-16\" src=\"/static/images/logo.svg\" alt=\"Logo\"> <span class=\"text-xl\">Go Feedback</span></a></div><div class=\"flex md:hidden\"><button class=\"btn btn-square btn-ghost\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = icons.Hamburger("").Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "</button></div><div class=\"hidden md:flex\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if auth.UseWaitList {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "<button onclick=\"auth_modal.showModal()\" class=\"transition-colors hover:border-none btn-block btn btn-neutral hover:bg-secondary hover:text-neutral\">Join Early Access</button>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else if auth.IsAuthenticated {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "<div class=\"flex-none\"><div class=\"dropdown dropdown-end\"><div tabindex=\"0\" role=\"button\" class=\"btn btn-ghost btn-circle avatar\"><div class=\"w-10 rounded-full\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if auth.AvatarURL != "" {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "<img alt=\"Profile Avatar\" src=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var2 string
				templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(auth.AvatarURL)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/transport/http/views/components/navbar.templ`, Line: 41, Col: 31}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 7, "<img alt=\"Profile Avatar\" src=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var3 string
				templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("https://api.dicebear.com/9.x/initials/svg?seed=%s", auth.Email))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/transport/http/views/components/navbar.templ`, Line: 46, Col: 93}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 8, "\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 9, "</div></div><ul tabindex=\"0\" class=\"p-2 mt-3 w-52 shadow menu menu-sm dropdown-content bg-base-100 rounded-box z-1\"><li><a class=\"justify-between\">Profile <span class=\"badge\">New</span></a></li><li><a>Settings</a></li><li><a href=\"/logout\">Logout</a></li></ul></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 10, "<div><a class=\"p-4 border-none btn btn-neutral hover:bg-secondary hover:text-neutral\" href=\"https://gofeedback-majiy.us.wristband.dev/login\">Log In</a> <a class=\"p-4 btn btn-outline\" href=\"https://gofeedback-majiy.us.wristband.dev/signup\">Sign Up</a></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 11, "</div></navbar>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = AuthModal().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 12, "</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
