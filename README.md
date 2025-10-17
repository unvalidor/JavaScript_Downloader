# JavaScript_Downloader

---

Download JavaScript files from a list of URLs — simple, fast, written in Go.

---

# Always check JavaScript files during recon

Whenever I'm working on a target, one of my fixed test-cases is to inspect the site's JavaScript files.
Why? Because libraries often contain vulnerable versions — finding a JS file that references an outdated dependency can quickly lead to a CVE and, in some cases, a real exploit.

I normally try to create a safe proof-of-concept (PoC) before reporting, but in practice many organizations are satisfied with a list of CVEs. That’s why including this test-case in our methodology is sensible and valuable.

---

# Step 1 — Collect URLs

Use the commands below to gather URLs from archives and live crawling:

```bash
# Wayback
echo "target.com" | waybackurls > wayback_urls.txt

# gau (getallurls)
echo target.com | gau --threads 10 --o gau.txt

# Katana for deeper crawling (depth 3, concurrency 50)
katana -u https://target.com -depth 3 -c 50 -o katana_urls.txt

# hakrawler for quick discovery and subdomains
echo "https://target.com" | hakrawler -d 3 -subs > hakrawler.txt
```

Now filter JavaScript files:

```bash
# Extract .js URLs from collected output
cat * | grep "\.js" > alljsfiles.txt

# Remove duplicates
cat alljsfiles.txt | sort -u > jsfiles.txt
```

---

# Step 2 — Download JavaScript files

Download all the JavaScript files you collected:

```bash
git clone https://github.com/unvalidor/JavaScript_Downloader
cd JavaScript_Downloader
go run js_downloader.go jsfiles.txt
```

Make sure the path to `jsfiles.txt` is correct.

---

# Final step — Scan for vulnerable libraries

Use `retire` to detect known vulnerable libraries in the downloaded files:

```bash
npm install -g retire
retire --version
retire --path ~/Desktop/recon/js_files
```

`~/Desktop/recon/js_files` should point to the folder where the downloader saved the JS files.

After the scan you will likely get a list of CVEs (if you're lucky). You can then search for PoCs (if you're lucky again 😄) or prepare reports.

**Important:** the mere presence of a CVE for a library version does **not** automatically mean the application is exploitable. Always triage results to determine runtime inclusion, exploitability, and business impact before escalating.

---

