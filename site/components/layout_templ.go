// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.546
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func navBar() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<header class=\"block w-1/2 mx-auto\"><div class=\"flex\"><a class=\"\" href=\"/recipe\">Recipes</a> <a class=\"pl-1\" href=\"/recipe/new\">Add Recipe</a></div></header>")
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
		Name: `__templ_onLoad_85a8`,
		Function: `function __templ_onLoad_85a8(){// htmx.logAll();

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
			if (name === "htmx:configRequest" ){
				console.log("json-enc-nested:",evt)
			}
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

	let src = new EventSource("/_templ/reload/events");
	src.onmessage = (event) => {
		if (event && event.data === "reload") {
			window.location.reload();
		}
	};
}`,
		Call:       templ.SafeScript(`__templ_onLoad_85a8`),
		CallInline: templ.SafeScriptInline(`__templ_onLoad_85a8`),
	}
}

func Layout() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var2 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var2 == nil {
			templ_7745c5c3_Var2 = templ.NopComponent
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
		var templ_7745c5c3_Var3 templ.ComponentScript = onLoad()
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var3.Call)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><script src=\"/assets/htmx.org@1.9.6\"></script><script src=\"/assets/hyperscript.org@0.9.12\"></script><script src=\"/assets/loading-states.js\"></script><script src=\"/assets/json-enc.js\"></script><script src=\"/assets/jquery-3.7.1.min.js\"></script><script src=\"/assets/serializeForm.js\"></script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = navBar().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"mx-auto w-3/4\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ_7745c5c3_Var2.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
