applicationId=$(curl https://discord.com/api/v10/applications/@me -H "Authorization: $DISCORD" | jq -r '.id')

curl -X PUT "https://discord.com/api/v10/applications/${applicationId}/commands" \
	-H "Authorization: $DISCORD" \
	-H "Content-Type: application/json" \
	-d '[
	{
		"name": "settings",
		"description": "Настройки"
	},
	{
		"name": "staff",
		"description": "Список учителей и персонала"
	},
	{
		"name": "stop",
		"description": "Остановить бота"
	},
	{
		"name": "students",
		"description": "Список учеников"
	}
]'
