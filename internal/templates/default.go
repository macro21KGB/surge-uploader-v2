package templates

func getDefaultTemplate() string {
	return `
<!doctype html>
<html>

<head>
	<title>Projects</title>
	<style>
		.card {
			border: 1px solid #ccc;
			padding: 10px;
			margin: 10px;
			border-radius: 5px;
		}

		#content {
			display: flex;
			flex-wrap: wrap;
			gap: 0.3rem;
		}

		.card-link {
			color: #a3a3ff;
			text-decoration: none;
		}
	</style>
</head>

<body>
	<template id="card">
		<div class="card">
			<p class="card-title">{{.Name}}</p>
			<a href="{{.URL}}" class="card-link">Visit</a>
		</div>
	</template>

	<div id="content"></div>
</body>

</html>
	`
}
