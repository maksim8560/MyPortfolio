package main

// Контент страницы живёт на сервере: Vue получает готовый JSON и только
// рендерит. Так текст правится в одном файле, без пересборки фронтенда.

type Link struct {
	Label    string `json:"label"`
	Href     string `json:"href"`
	External bool   `json:"external,omitempty"`
}

type Fact struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Stat struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type Project struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Kind     string   `json:"kind"`
	Tagline  string   `json:"tagline"`
	Summary  string   `json:"summary"`
	Bullets  []string `json:"bullets"`
	Tags     []string `json:"tags"`
	Demo     string   `json:"demo,omitempty"`
	DemoText string   `json:"demoText,omitempty"`
	Year     string   `json:"year"`
	Status   string   `json:"status"`
	Featured bool     `json:"featured"`
}

type Approach struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type ContactItem struct {
	Label  string `json:"label"`
	Value  string `json:"value"`
	Href   string `json:"href,omitempty"`
	Copy   bool   `json:"copy,omitempty"`
	Copied string `json:"copied,omitempty"`
}

type Section struct {
	Kicker string `json:"kicker"`
	Title  string `json:"title"`
	Lead   string `json:"lead"`
}

type Content struct {
	Lang     string          `json:"lang"`
	Meta     Meta            `json:"meta"`
	Nav      []Link          `json:"nav"`
	Hero     Hero            `json:"hero"`
	Stats    []Stat          `json:"stats"`
	Work     Section         `json:"work"`
	Projects []Project       `json:"projects"`
	Approach ApproachSection `json:"approach"`
	Contact  ContactSection  `json:"contact"`
	Footer   Footer          `json:"footer"`
}

// ApproachSection и ContactSection вынесены в именованные типы,
// чтобы админка могла работать с ними как с обычными полями.
type ApproachSection struct {
	Section
	Items []Approach `json:"items"`
}

type ContactSection struct {
	Section
	Items []ContactItem `json:"items"`
	Note  string        `json:"note"`
}

type Meta struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Hero struct {
	Kicker    string `json:"kicker"`
	Title1    string `json:"title1"`
	Title2    string `json:"title2"`
	Lead      string `json:"lead"`
	Primary   Link   `json:"primary"`
	Secondary Link   `json:"secondary"`
	Facts     []Fact `json:"facts"`
	Note      string `json:"note"`
}

type Footer struct {
	Note  string `json:"note"`
	Built string `json:"built"`
}

// defaultRU и defaultEN — стартовый контент. Используется только при первом
// запуске: дальше правки живут в data/content.json и не перетираются.
func defaultRU() Content {
	c := Content{
		Lang: "ru",
		Meta: Meta{
			Title:       "MaxDiWay (Максим) — инди-разработчик",
			Description: "Портфолио: Aura — мессенджер с E2E-шифрованием и WebRTC-звонками, Sonora и FileDropper. Vue на фронте, Go на бэке.",
		},
		Nav: []Link{
			{Label: "Проекты", Href: "#work"},
			{Label: "Подход", Href: "#approach"},
			{Label: "Контакты", Href: "#contact"},
		},
		Hero: Hero{
			Kicker:  "Обо мне",
			Title1:  "Инди-разработчик",
			Title2:  "полного цикла",
			Lead:    "Шифрование, которое действительно шифрует. Загрузки, которые не упираются в лимит платформы. Интерфейс, который не притворяется, что знает больше, чем знает.",
			Primary: Link{Label: "Смотреть проекты", Href: "#work"},
			Secondary: Link{
				Label:    "Открыть мессенджер",
				Href:     "https://aura-messenger-q5yk.onrender.com/",
				External: true,
			},
			Facts: []Fact{
				{Key: "фронтенд", Value: "Vue 3 · TypeScript · Vite"},
				{Key: "бэкенд", Value: "Go · net/http · embed"},
				{Key: "статус", Value: "Aura в публичном проде"},
			},
			Note: "Aura — мессенджер, который можно открыть и потрогать прямо сейчас.",
		},
		Stats: []Stat{
			{Value: "3", Label: "проекта в публичном проде"},
			{Value: "E2E", Label: "шифрование в мессенджере"},
			{Value: "0", Label: "npm-зависимостей у Sonora"},
			{Value: "2", Label: "языка интерфейса"},
		},
		Work: Section{
			Kicker: "Проекты",
			Title:  "Что сделано",
			Lead:   "Три продукта, каждый доведён до публичного деплоя. Ниже — не список фич, а то, что на самом деле оказалось сложным.",
		},
		Projects: []Project{
			{
				ID:      "aura",
				Name:    "Aura",
				Kind:    "Мессенджер",
				Tagline: "Мессенджер с E2E-шифрованием, WebRTC-звонками и голосовыми заметками в liquid-glass интерфейсе.",
				Summary: "Мой основной проект. Личные и групповые чаты, звонки один-на-один и со скриншерингом, голосовые, вложения, роли и модерация. Шифрование на клиенте, ключи не покидают устройство, а сервер видит только шифротекст.",
				Bullets: []string{
					"Сквозное шифрование на WebCrypto: сервер хранит шифротекст и не может прочитать переписку.",
					"Звонки на WebRTC с собственным STUN/TURN-клиентом и диагностикой ICE: «Allocate 401» — это запрос пароля, а не сломанный сервер.",
					"«Отправлено» значит «сохранено»: ack ждёт записи на диск и в зеркало, тест убивает процесс SIGKILL внутри окна записи.",
					"Мобильный клиент и PWA, безопасность сессий, ограничение частоты запросов и защита от двойной отправки.",
				},
				Tags:     []string{"WebRTC", "WebCrypto", "Node.js", "PWA", "Render"},
				Demo:     "https://aura-messenger-q5yk.onrender.com/",
				DemoText: "Открыть Aura",
				Year:     "2025—2026",
				Status:   "в проде",
				Featured: true,
			},
			{
				ID:      "sonora",
				Name:    "Sonora",
				Kind:    "Музыкальный плеер",
				Tagline: "Плеер без npm-зависимостей: свой бандлер, свой парсер ID3, свой аудио-граф.",
				Summary: "Эксперимент на тему «сколько можно сделать, не подключая чужие пакеты». Спектрограмма рисуется только у потока, который реально прошёл через аудио-граф, — иначе её просто нет.",
				Bullets: []string{
					"19 CSS-токенов на @property вместо рантайм-темы.",
					"Свой парсер ID3 и обложек — без сторонних библиотек.",
					"Честные состояния: если данных ещё нет, интерфейс не рисует ноль.",
				},
				Tags:     []string{"Vanilla JS", "Web Audio", "@property", "GitHub Pages"},
				Demo:     "https://maksim8560.github.io/music-app/",
				DemoText: "Открыть Sonora",
				Year:     "2025",
				Status:   "в проде",
			},
			{
				ID:      "filedropper",
				Name:    "FileDropper",
				Kind:    "Обмен файлами",
				Tagline: "Передача файлов через Cloudflare Workers, где байты не проходят через воркер.",
				Summary: "Загрузка идёт напрямую в хранилище по подписанным ссылкам, воркер только выдаёт разрешение и следит за лимитами. Поэтому файл не упирается в лимит тела запроса платформы.",
				Bullets: []string{
					"Presigned URL: воркер не проксирует содержимое файла.",
					"Лимиты и срок жизни ссылки считаются на сервере.",
					"SHA-адресация и аккуратная уборка просроченного.",
				},
				Tags:     []string{"Cloudflare Workers", "R2", "Presigned URL", "Vue"},
				Demo:     "https://filedropper-api.sonora-online.workers.dev/#/",
				DemoText: "Открыть FileDropper",
				Year:     "2025",
				Status:   "в проде",
			},
		},
		Contact: ContactSection{
			Section: Section{
				Kicker: "Контакты",
				Title:  "Написать",
				Lead:   "Открыт к задачам, где нужно довести продукт до продакшена целиком: от схемы данных до последнего пикселя.",
			},
			Items: []ContactItem{
				{Label: "Discord", Value: "MaxDiWay", Copy: true, Copied: "Скопировано"},
				{Label: "GitHub", Value: "maksim8560", Href: "https://github.com/maksim8560", Copy: false},
			},
			Note: "Discord — самый быстрый способ: ник копируется в один клик.",
		},
		Footer: Footer{
			Note:  "Собрано на Vue 3 и Go. Контент отдаёт Go, статику — embed.",
			Built: "Vue + Go",
		},
	}
	c.Approach.Section = Section{
		Kicker: "Подход",
		Title:  "Принципы",
		Lead:   "Каждый пункт вырос из конкретного отказа, а не из методологии: это то, что пришлось переделать, потому что «и так работало» оказалось неправдой.",
	}
	c.Approach.Items = []Approach{
		{Title: "«Настроено» ≠ «работает»", Body: "Зеркало, TURN, хранилище, почта — каждая подсистема отдаёт configured и reachable отдельными полями. Однажды хранилище отвечало 502 целый день, а код проверял только, что переменные окружения проставлены."},
		{Title: "Долговечность важнее удобства API", Body: "Ответ клиенту ждёт, пока запись действительно легла на диск и в зеркало. Тест растягивает окно до трёх секунд, убивает процесс SIGKILL внутри него и требует, чтобы переписка вернулась целиком."},
		{Title: "Честный интерфейс лучше красивого", Body: "Счётчик не показывает ноль, пока сервер не ответил «нет данных». Плеер не рисует спектрограмму у потока, который не прошёл через аудио-граф. Если система деградировала — на странице табличка с командой, а не пустая 500."},
		{Title: "Зависимость — это решение", Body: "Один из проектов живёт вообще без npm-пакетов: свой бандлер, свой парсер ID3, свой STUN/TURN-клиент. Не из любви к аскетизму — лишняя зависимость в чужом реестре обновляется не в твой день релизов."},
		{Title: "Тест, который не может упасть, — не тест", Body: "В Aura больше тридцати автоматических проверок: двойная отправка, порядок сообщений, истечение сессии, права в группах, восстановление после падения. Проверка, которая проходит всегда, удаляется."},
		{Title: "Сервер — единственный источник правды", Body: "Клиентский UI не решает, можно ли действие: он показывает то, что разрешил сервер. Это и защита от подделки запросов, и способ не расходиться двум вкладкам в состоянии."},
	}
	return c
}

func defaultEN() Content {
	c := Content{
		Lang: "en",
		Meta: Meta{
			Title:       "MaxDiWay (Maxim) — indie developer",
			Description: "Portfolio: Aura, an end-to-end encrypted messenger with WebRTC calls, plus Sonora and FileDropper. Vue on the front end, Go on the back end.",
		},
		Nav: []Link{
			{Label: "Work", Href: "#work"},
			{Label: "Approach", Href: "#approach"},
			{Label: "Contact", Href: "#contact"},
		},
		Hero: Hero{
			Kicker:  "About",
			Title1:  "Indie developer",
			Title2:  "end to end",
			Lead:    "Encryption that actually encrypts. Uploads that do not hit the platform's request limit. An interface that does not pretend to know more than it does.",
			Primary: Link{Label: "See the work", Href: "#work"},
			Secondary: Link{
				Label:    "Open the messenger",
				Href:     "https://aura-messenger-q5yk.onrender.com/",
				External: true,
			},
			Facts: []Fact{
				{Key: "frontend", Value: "Vue 3 · TypeScript · Vite"},
				{Key: "backend", Value: "Go · net/http · embed"},
				{Key: "status", Value: "Aura live in production"},
			},
			Note: "Aura is a messenger you can open and poke at right now.",
		},
		Stats: []Stat{
			{Value: "3", Label: "projects in public production"},
			{Value: "E2E", Label: "encryption in the messenger"},
			{Value: "0", Label: "npm dependencies in Sonora"},
			{Value: "2", Label: "interface languages"},
		},
		Work: Section{
			Kicker: "Work",
			Title:  "What I built",
			Lead:   "Three products, each shipped to a public deploy. Below is not a feature list — it is what actually turned out to be hard.",
		},
		Projects: []Project{
			{
				ID:      "aura",
				Name:    "Aura",
				Kind:    "Messenger",
				Tagline: "An end-to-end encrypted messenger with WebRTC calls and voice notes behind a liquid-glass interface.",
				Summary: "My main project. Private and group chats, one-to-one and screen-sharing calls, voice notes, attachments, roles and moderation. Encryption happens on the client; keys never leave the device and the server only ever sees ciphertext.",
				Bullets: []string{
					"End-to-end encryption on WebCrypto: the server stores ciphertext and cannot read a conversation.",
					"Calls on WebRTC with a hand-written STUN/TURN client and ICE diagnostics: \"Allocate 401\" is a password request, not a broken server.",
					"\"Sent\" means \"stored\": the ack waits for the disk write and the mirror, and a test SIGKILLs the process inside that window.",
					"Mobile client and PWA, session security, rate limiting and duplicate-send protection.",
				},
				Tags:     []string{"WebRTC", "WebCrypto", "Node.js", "PWA", "Render"},
				Demo:     "https://aura-messenger-q5yk.onrender.com/",
				DemoText: "Open Aura",
				Year:     "2025—2026",
				Status:   "live",
				Featured: true,
			},
			{
				ID:      "sonora",
				Name:    "Sonora",
				Kind:    "Music player",
				Tagline: "A player with no npm dependencies: its own bundler, its own ID3 parser, its own audio graph.",
				Summary: "An experiment about how far you can get without third-party packages. The spectrogram only appears for a stream that actually passed through the audio graph — otherwise there simply is none.",
				Bullets: []string{
					"19 CSS tokens on @property instead of a runtime theme.",
					"Hand-written ID3 and cover-art parsing — no libraries.",
					"Honest states: with no data yet, the interface does not draw a zero.",
				},
				Tags:     []string{"Vanilla JS", "Web Audio", "@property", "GitHub Pages"},
				Demo:     "https://maksim8560.github.io/music-app/",
				DemoText: "Open Sonora",
				Year:     "2025",
				Status:   "live",
			},
			{
				ID:      "filedropper",
				Name:    "FileDropper",
				Kind:    "File transfer",
				Tagline: "File transfer on Cloudflare Workers where the bytes never pass through the worker.",
				Summary: "Uploads go straight to storage over presigned URLs; the worker only issues the permission and tracks the limits. That is why a file never hits the platform's request body limit.",
				Bullets: []string{
					"Presigned URLs: the worker does not proxy file contents.",
					"Quotas and link lifetime are enforced server-side.",
					"Content addressing and tidy expiry cleanup.",
				},
				Tags:     []string{"Cloudflare Workers", "R2", "Presigned URL", "Vue"},
				Demo:     "https://filedropper-api.sonora-online.workers.dev/#/",
				DemoText: "Open FileDropper",
				Year:     "2025",
				Status:   "live",
			},
		},
		Contact: ContactSection{
			Section: Section{
				Kicker: "Contact",
				Title:  "Get in touch",
				Lead:   "Open to work where a product has to be taken all the way to production: from the data model to the last pixel.",
			},
			Items: []ContactItem{
				{Label: "Discord", Value: "MaxDiWay", Copy: true, Copied: "Copied"},
				{Label: "GitHub", Value: "maksim8560", Href: "https://github.com/maksim8560"},
			},
			Note: "Discord is the fastest route: the handle copies in one click.",
		},
		Footer: Footer{
			Note:  "Built with Vue 3 and Go. Go serves the content, embed serves the static files.",
			Built: "Vue + Go",
		},
	}
	c.Approach.Section = Section{
		Kicker: "Approach",
		Title:  "Principles",
		Lead:   "Every item below grew out of a specific failure, not a methodology: it is what had to be rebuilt because \"it already worked\" turned out to be untrue.",
	}
	c.Approach.Items = []Approach{
		{Title: "\"Configured\" is not \"working\"", Body: "Mirror, TURN, storage and mail each report configured and reachable as separate fields. Once storage returned 502 for a whole day while the code only checked that environment variables were set."},
		{Title: "Durability beats API convenience", Body: "The response waits until the write has really landed on disk and in the mirror. A test stretches that window to three seconds, SIGKILLs the process inside it and requires the conversation to come back whole."},
		{Title: "An honest interface beats a pretty one", Body: "A counter does not show zero until the server answered \"no data\". A player does not draw a spectrogram for a stream that never passed through the audio graph. When the system degrades, the page shows the command to run instead of a blank 500."},
		{Title: "A dependency is a decision", Body: "One project runs with no npm packages at all: its own bundler, its own ID3 parser, its own STUN/TURN client. Not out of asceticism — a dependency in someone else's registry does not update on your release day."},
		{Title: "A test that cannot fail is not a test", Body: "Aura has over thirty automated checks: double sends, message ordering, session expiry, group permissions, crash recovery. A check that always passes gets deleted."},
		{Title: "The server is the only source of truth", Body: "The client UI does not decide whether an action is allowed; it shows what the server permitted. That is both forged-request protection and the reason two tabs never disagree."},
	}
	return c
}
