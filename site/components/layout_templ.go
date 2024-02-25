// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.571
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "github.com/JamesTiberiusKirk/recipe-cms/common"
import "github.com/JamesTiberiusKirk/recipe-cms/site/session"
import "fmt"

func navBar(c *common.Context) templ.Component {
	sess, ok := c.Get("session").(*session.Manager)
	if !ok {
		fmt.Println("not ok")
		return navBarTempl(false, "")
	}

	user, err := sess.GetUser(c)
	if err != nil {
		fmt.Println("error getting user")
		return navBarTempl(false, "")
	}

	if user == "" {
		fmt.Println("user is empty ")
		return navBarTempl(false, "")
	}

	return navBarTempl(true, user)
}

func navBarTempl(loggedIn bool, user string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"w-full\"><header class=\"lg:w-1/2 sm:w-full mx-auto\"><div class=\"lg:text-xl text-5xl flex p-2 justify-evenly\"><a class=\"p-6\" href=\"/recipes\">Recipes</a> <a class=\"p-6\" href=\"/recipe/new\">Add Recipe</a> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if loggedIn {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<a class=\"p-6\" href=\"/auth/logout\">Logout</a> <span class=\"p-6 text-sm\">USER: ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var2 string
			templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(user)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `site/components/layout.templ`, Line: 35, Col: 43}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<a class=\"p-6\" href=\"/auth/login\">Login</a>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></header></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func onLoad() templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_onLoad_4ce8`,
		Function: `function __templ_onLoad_4ce8(){// htmx.logAll();

	htmx.config.defaultSettleDelay = 0;
	//document.addEventListener('htmx:afterRequest', function (event) {
	// console.log(event)
	//})
	document.addEventListener('htmx:beforeOnLoad', function (event) {
			// console.log("htmx:beforeOnLoad CUSTOM", event)

			const res = event.detail.xhr
			const isValidationRoute = (event.detail.requestConfig.path.search("validate"))

			if (isValidationRoute) {
				if (res.status === 400) {
					// console.log("status === 400")
					event.detail.shouldSwap = true;
					event.detail.isError = false;
				}

				if (res.status === 200) {
					// console.log("status === 200")
				}
			}

			// Bellow is an example of creating your own custom attibutes
			// In this example I have created attibutes which dynamically enables or disables (buttons) when using validation routes
			// This same could be expanded for a lot of other behaviour
			// However, this should not be overused as it could already be covered in htmx
			// AND because too much js here will just render the whole point of htmx nil
			if (isValidationRoute) {
				if (res.status === 400) {
					const disableOnErrorTargetId = event.srcElement.getAttribute("disableOnError")
					if (disableOnErrorTargetId) {
						const targetIdParsed = disableOnErrorTargetId.replace("#", "")
							const target = document.getElementById(targetIdParsed)
							if (target) target.setAttribute("disabled", true)
					}
				}

				if (res.status === 200) {
				const enableOnValidTargetId = event.srcElement.getAttribute("enableOnValid")
					if (enableOnValidTargetId) {
						const targetIdParsed = enableOnValidTargetId.replace("#", "")
						const target = document.getElementById(targetIdParsed)
						if (target) target.removeAttribute("disabled")
					}
				}
			}
	});

	htmx.defineExtension('markdown-preview', {
		onEvent: function (name, evt) {
			if (name === "htmx:configRequest"){
				console.log("markdown configRequest",evt)
				evt.detail.headers['Content-Type'] = "text/plain";
			}
		},
		encodeParameters : function(xhr, parameters, elt) {
			console.log("markdown elt.value:", elt.value)
			if (elt.tagName !== "FORM"){
				xhr.overrideMimeType('text/plain');
				return encodeURIComponent(elt.value)
			}
		}
	});

	htmx.defineExtension('json-enc-nested', {
		onEvent: function (name, evt) {
			if (name === "htmx:configRequest" &&
				(evt.srcElement.tagName === "FORM" ||
				 evt.srcElement.form != undefined) ) {
				console.log("json-enc-nested:",evt)
				evt.detail.headers['Content-Type'] = "application/json";
			}
		},
		encodeParameters : function(xhr, parameters, elt) {
			console.log("json-enc-nested elt:")
			console.dir(elt)

			if (elt.form === undefined) return

			let jsonEnc = $(elt.form).serializeJSON()
			console.log("json-enc-nested",jsonEnc)

			xhr.overrideMimeType('application/json');
			return JSON.stringify(jsonEnc);
		}
	});

	// NOTE: not actually intended to have multiple at once 
	// Might need to completly set the toast from js if i want it to work like a stack 
	const toaster = (e) => {showToast(e.type, e.detail.value)} 
	htmx.on("toast_danger", toaster)
	htmx.on("toast_warning", toaster)
	htmx.on("toast_success", toaster)
	htmx.on("toast_info", toaster)

        function showToast(type, message) {
		el = document.getElementById(type)
		document.getElementById(type+"_text").innerHTML = message
		el.classList.remove("hidden");
		setTimeout(function () {
			el.classList.add("hidden");
			document.getElementById(type+"-text").innerHTML = ""
		}, 5000);
        }
}`,
		Call:       templ.SafeScript(`__templ_onLoad_4ce8`),
		CallInline: templ.SafeScriptInline(`__templ_onLoad_4ce8`),
	}
}

func dismissToast(id string) templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_dismissToast_7770`,
		Function: `function __templ_dismissToast_7770(id){document.getElementById(id).classList.add("hidden");
	document.getElementById(id+"_text").innerHTML = "";
}`,
		Call:       templ.SafeScript(`__templ_dismissToast_7770`, id),
		CallInline: templ.SafeScriptInline(`__templ_dismissToast_7770`, id),
	}
}

func toasts() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var3 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var3 == nil {
			templ_7745c5c3_Var3 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"toast_danger\" class=\"hidden z-10 absolute bottom-0 left-1/2 -translate-x-1/2 my-auto\"><div class=\"mb-4 bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative\" role=\"alert\"><span id=\"toast_danger_text\" class=\"block sm:inline pr-6\"></span> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ.RenderScriptItems(ctx, templ_7745c5c3_Buffer, dismissToast("toast-danger"))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span onclick=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 templ.ComponentScript = dismissToast("toast-danger")
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var4.Call)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" class=\"absolute top-0 bottom-0 right-0 px-4 py-3\"><svg class=\"fill-current h-6 w-6 text-red-500\" role=\"button\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 20 20\"><title>Close</title><path d=\"M14.348 14.849a1.2 1.2 0 0 1-1.697 0L10 11.819l-2.651 3.029a1.2 1.2 0 1 1-1.697-1.697l2.758-3.15-2.759-3.152a1.2 1.2 0 1 1 1.697-1.697L10 8.183l2.651-3.031a1.2 1.2 0 1 1 1.697 1.697l-2.758 3.152 2.758 3.15a1.2 1.2 0 0 1 0 1.698z\"></path></svg></span></div></div><div id=\"toast_warning\" class=\"hidden z-10 absolute bottom-0 left-1/2 -translate-x-1/2 my-auto\"><div class=\"mb-4 bg-yellow-100 border border-yellow-400 text-yellow-700 px-4 py-3 rounded relative\" role=\"alert\"><span id=\"toast_warning_text\" class=\"block sm:inline pr-6\"></span> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ.RenderScriptItems(ctx, templ_7745c5c3_Buffer, dismissToast("toast-warning"))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span onclick=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var5 templ.ComponentScript = dismissToast("toast-warning")
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var5.Call)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" class=\"absolute top-0 bottom-0 right-0 px-4 py-3\"><svg class=\"fill-current h-6 w-6 text-yellow-500\" role=\"button\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 20 20\"><title>Close</title><path d=\"M14.348 14.849a1.2 1.2 0 0 1-1.697 0L10 11.819l-2.651 3.029a1.2 1.2 0 1 1-1.697-1.697l2.758-3.15-2.759-3.152a1.2 1.2 0 1 1 1.697-1.697L10 8.183l2.651-3.031a1.2 1.2 0 1 1 1.697 1.697l-2.758 3.152 2.758 3.15a1.2 1.2 0 0 1 0 1.698z\"></path></svg></span></div></div><div id=\"toast_success\" class=\"hidden z-10 absolute bottom-0 left-1/2 -translate-x-1/2 my-auto\"><div class=\"mb-4 bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded relative\" role=\"alert\"><span id=\"toast_success_text\" class=\"block sm:inline pr-6\"></span> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ.RenderScriptItems(ctx, templ_7745c5c3_Buffer, dismissToast("toast-success"))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span onclick=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var6 templ.ComponentScript = dismissToast("toast-success")
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var6.Call)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" class=\"absolute top-0 bottom-0 right-0 px-4 py-3\"><svg class=\"fill-current h-6 w-6 text-green-500\" role=\"button\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 20 20\"><title>Close</title><path d=\"M14.348 14.849a1.2 1.2 0 0 1-1.697 0L10 11.819l-2.651 3.029a1.2 1.2 0 1 1-1.697-1.697l2.758-3.15-2.759-3.152a1.2 1.2 0 1 1 1.697-1.697L10 8.183l2.651-3.031a1.2 1.2 0 1 1 1.697 1.697l-2.758 3.152 2.758 3.15a1.2 1.2 0 0 1 0 1.698z\"></path></svg></span></div></div><div id=\"toast_info\" class=\"hidden z-10 absolute bottom-0 left-1/2 -translate-x-1/2 my-auto\"><div class=\"mb-4 bg-blue-100 border border-blue-400 text-blue-700 px-4 py-3 rounded relative\" role=\"alert\"><span id=\"toast_info_text\" class=\"block sm:inline pr-6\"></span> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ.RenderScriptItems(ctx, templ_7745c5c3_Buffer, dismissToast("toast-info"))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<span onclick=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var7 templ.ComponentScript = dismissToast("toast-info")
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var7.Call)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" class=\"absolute top-0 bottom-0 right-0 px-4 py-3\"><svg class=\"fill-current h-6 w-6 text-blue-500\" role=\"button\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 20 20\"><title>Close</title><path d=\"M14.348 14.849a1.2 1.2 0 0 1-1.697 0L10 11.819l-2.651 3.029a1.2 1.2 0 1 1-1.697-1.697l2.758-3.15-2.759-3.152a1.2 1.2 0 1 1 1.697-1.697L10 8.183l2.651-3.031a1.2 1.2 0 1 1 1.697 1.697l-2.758 3.152 2.758 3.15a1.2 1.2 0 0 1 0 1.698z\"></path></svg></span></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func Layout(c *common.Context) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var8 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var8 == nil {
			templ_7745c5c3_Var8 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<html class=\"dark\"><head><link rel=\"stylesheet\" href=\"/assets/styles.css\"><link rel=\"stylesheet\" href=\"/assets/tailwind.css\"></head>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ.RenderScriptItems(ctx, templ_7745c5c3_Buffer, onLoad())
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<body class=\"dark:bg-gray-800 dark:text-white\" onload=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var9 templ.ComponentScript = onLoad()
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var9.Call)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><script src=\"/assets/htmx.org@1.9.6\"></script><script src=\"/assets/hyperscript.org@0.9.12\"></script><script src=\"/assets/loading-states.js\"></script><script src=\"/assets/json-enc.js\"></script><script src=\"/assets/jquery-3.7.1.min.js\"></script><script src=\"/assets/serializeForm.js\"></script><script src=\"/assets/Sortable.js\"></script><script src=\"https://unpkg.com/htmx.org/dist/ext/sse.js\"></script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = navBar(c).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"mx-auto w-11/12 lg:w-3/4 sm:text-xl\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ_7745c5c3_Var8.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = toasts().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
