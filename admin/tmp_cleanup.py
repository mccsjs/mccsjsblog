import re

with open('src/views/setting/Setting.vue', 'r', encoding='utf-8') as f:
    content = f.read()

# 1. Remove mottoMainList leftover lines (.fill(null) and .map line)
# These are orphan lines left after removing the mottoMainList assignment
content = re.sub(
    r'\n\s+\.fill\(null\)\s*\n\s+\.map\(\(_\s*,\s*i\)\s*=>\s*parseJSON\(configs\.about_motto_main.*?\);',
    '',
    content,
    flags=re.DOTALL
)

# Also remove the orphan blank line before socializeList if present
content = re.sub(
    r'(\.footerLinksList.*?\]\);\s*\n)\n(\s*blogForm\.value\.socializeList)',
    r'\1\2',
    content,
    flags=re.DOTALL
)

# 2. Remove aboutExhibitionUploaderRef upload promise (lines ~538-542)
content = re.sub(
    r',\s*\n\s*uploadPromises\.push\(\s*\n\s*blogUploaders\.aboutExhibitionUploaderRef\.uploadPendingFile\(\)\.then\(url\s*=>\s*\{[^}]*\}\)\s*\);',
    '',
    content,
    flags=re.DOTALL
)

# 3. Remove save payload fields
payload_removals = [
    r"'blog\.about_exhibition':\s*blogForm\.value\.about_exhibition,\s*\n",
    r"'blog\.about_profile':\s*JSON\.stringify\(blogForm\.value\.profileList\),\s*\n",
    r"'blog\.about_personality':\s*blogForm\.value\.about_personality,\s*\n",
    r"'blog\.about_motto_main':\s*JSON\.stringify\(blogForm\.value\.mottoMainList\),\s*\n",
    r"'blog\.about_motto_sub':\s*blogForm\.value\.about_motto_sub,\s*\n",
    r"'blog\.about_unions':\s*JSON\.stringify\(blogForm\.value\.unionsList\),\s*\n",
]

for pattern in payload_removals:
    content = re.sub(pattern, '', content)

with open('src/views/setting/Setting.vue', 'w', encoding='utf-8') as f:
    f.write(content)

print('All cleanup done!')
print('Checking for remaining references...')

# Verify
remaining = re.findall(r'about_exhibition|about_personality|about_motto|about_profile|about_unions|profileList|mottoMainList|unionsList', content)
if remaining:
    print(f'WARNING: Still found {len(remaining)} references:')
    for r in set(remaining):
        print(f'  - {r}')
else:
    print('No remaining references found - clean!')
