// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.833

package icons

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import (
	"github.com/a-h/templ"
	templruntime "github.com/a-h/templ/runtime"
)

func AIBrain(className string) templ.Component {
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
		var templ_7745c5c3_Var2 = []any{className}
		templ_7745c5c3_Err = templ.RenderCSSItems(ctx, templ_7745c5c3_Buffer, templ_7745c5c3_Var2...)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 24 24\" fill=\"none\" class=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(templ.CSSClasses(templ_7745c5c3_Var2).String())
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/transport/http/views/icons/aibrain.templ`, Line: 1, Col: 0}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "\"><path fill-rule=\"evenodd\" clip-rule=\"evenodd\" d=\"M9.5 22.7499C10.5052 22.7499 11.4039 22.2935 12 21.5766C12.5962 22.2934 13.4948 22.7498 14.5 22.7498C16.0585 22.7498 17.3608 21.6528 17.6768 20.1888C19.4248 19.8702 20.75 18.3397 20.75 16.4998C20.75 16.1171 20.6925 15.7473 20.5856 15.3988C21.8641 14.8014 22.75 13.5046 22.75 11.9998C22.75 10.4949 21.8641 9.19809 20.5856 8.60071C20.6925 8.25221 20.75 7.88237 20.75 7.49976C20.75 5.65977 19.4248 4.12929 17.6768 3.81067C17.3608 2.34673 16.0585 1.24976 14.5 1.24976C13.4948 1.24976 12.5961 1.70617 12 2.42305C11.4038 1.70621 10.5052 1.24985 9.5 1.24985C7.94152 1.24985 6.63925 2.34683 6.3232 3.81077C4.57518 4.12938 3.25 5.65986 3.25 7.49985C3.25 7.88246 3.30751 8.2523 3.41442 8.6008C2.13588 9.19818 1.25 10.495 1.25 11.9999C1.25 13.5047 2.13588 14.8015 3.41442 15.3989C3.30751 15.7474 3.25 16.1172 3.25 16.4999C3.25 18.3398 4.57518 19.8703 6.3232 20.1889C6.63925 21.6529 7.94152 22.7499 9.5 22.7499ZM10.0002 7.74976C9.37879 7.74976 8.82709 8.1474 8.63058 8.73693L6.78869 14.2626C6.65771 14.6555 6.87008 15.0803 7.26303 15.2113C7.65599 15.3423 8.08073 15.1299 8.21172 14.7369L8.70744 13.2498H11.293L11.7887 14.7369C11.9197 15.1299 12.3444 15.3423 12.7374 15.2113C13.1303 15.0803 13.3427 14.6555 13.2117 14.2626L11.3698 8.73693C11.1733 8.1474 10.6216 7.74976 10.0002 7.74976ZM10.0002 9.37146L10.793 11.7498H9.20744L10.0002 9.37146ZM16.2502 8.49976C16.2502 8.08554 15.9144 7.74976 15.5002 7.74976C15.086 7.74976 14.7502 8.08554 14.7502 8.49976V14.4998C14.7502 14.914 15.086 15.2498 15.5002 15.2498C15.9144 15.2498 16.2502 14.914 16.2502 14.4998V8.49976Z\" fill=\"currentColor\"></path></svg>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
