---
layout: page
title: "Osprey video tutorials: watch real apps get built"
description: Video tutorials for the Osprey language. Watch real Osprey apps get built, see the code on screen, then run each one yourself with the commands shown.
permalink: /docs/tutorials/
tags: [tutorials, video]
---

Each tutorial is a short narrated video built around one real Osprey app. The video shows the app running and quotes the code from the repository as it goes. Under each video, the page lists the chapters, the commands to run the app yourself, and a full transcript.

{%- for item in collections.all | sort(false, false, "data.title") %}{% if item.data.video and item.url.startsWith("/docs/tutorials/") %}{% set v = item.data.video | video %}

## {{ v.name }}

{% videoCard item.data.video, item.url %}

{{ v.description }}
{%- endif %}{% endfor %}

## More to read

The [documentation home](/docs/) lists the written guides, including [web apps](/docs/web-apps/), [iOS and Android apps](/docs/mobile-apps/) and [algebraic effects](/docs/effects/).
