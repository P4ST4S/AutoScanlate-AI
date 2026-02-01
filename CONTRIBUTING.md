# 🤝 Contributing to Manga AI Translator

Merci de votre intérêt pour contribuer au projet ! / Thank you for your interest in contributing to the project!

_Ce document est disponible en français et en anglais. / This document is available in French and English._

---

## 🌍 Languages / Langues

Vous pouvez communiquer en **français** ou en **anglais** dans vos issues, PRs et commentaires.

You can communicate in **French** or **English** in your issues, PRs, and comments.

---

## 📋 Comment contribuer / How to Contribute

### 1. 💡 Proposer une modification / Proposing a Change

**Il est fortement recommandé d'ouvrir une Issue avant de commencer à coder**, surtout pour les fonctionnalités importantes. Cela évite de perdre du temps sur quelque chose qui pourrait ne pas être accepté.

**It is strongly recommended to open an Issue before starting to code**, especially for significant features. This avoids wasting time on something that might not be accepted.

**⚠️ Zone sensible / Sensitive Area** : Le module `/ai-worker` est la partie la plus critique du projet. Toute modification dans ce module **nécessite une discussion préalable** via une Issue.

**⚠️ Sensitive Area**: The `/ai-worker` module is the most critical part of the project. Any changes to this module **require prior discussion** via an Issue.

### 2. 🔧 Types de contributions acceptées / Accepted Contribution Types

Nous acceptons tous les types de contributions si elles sont **bien justifiées et présentées** :

We accept all types of contributions if they are **well-justified and presented**:

- 🐛 **Bug fixes** : Corrections de bugs / Bug corrections
- ✨ **Features** : Nouvelles fonctionnalités / New features
- 📚 **Documentation** : Améliorations de la documentation / Documentation improvements
- 🌐 **Translations** : Traductions / Translations
- ⚡ **Performance** : Optimisations / Optimizations
- 🧪 **Tests** : Ajout ou amélioration des tests / Adding or improving tests
- 🎨 **UI/UX** : Améliorations de l'interface / Interface improvements

---

## 🛠️ Standards techniques / Technical Standards

### Code Style

Nous utilisons les outils suivants pour maintenir la qualité du code :

We use the following tools to maintain code quality:

- **Python** (`/ai-worker`) : [Ruff](https://docs.astral.sh/ruff/) pour le linting et le formatage
- **Go** (`/backend-api`) : [golangci-lint](https://golangci-lint.run/) et `gofmt`
- **TypeScript/React** (`/frontend`) : [ESLint](https://eslint.org/)

**Assurez-vous que votre code passe les linters avant de soumettre une PR.**

**Make sure your code passes the linters before submitting a PR.**

```bash
# Python
cd ai-worker
ruff check .
ruff format .

# Go
cd backend-api
golangci-lint run
gofmt -w .

# Frontend
cd frontend
pnpm lint
```

### Tests

**Les nouvelles fonctionnalités et bug fixes doivent inclure des tests.**

**New features and bug fixes must include tests.**

- **Python** : Utilisez `pytest`
- **Go** : Utilisez le framework de test standard Go (`testing`)
- **Frontend** : Utilisez Jest ou React Testing Library (à configurer si nécessaire)

### Documentation

**Toute nouvelle fonctionnalité doit être documentée** :

**All new features must be documented**:

- Ajoutez des commentaires dans le code pour les parties complexes
- Mettez à jour le README approprié si nécessaire
- Documentez les nouvelles APIs ou endpoints

---

## 📝 Processus de Pull Request / Pull Request Process

### 1. Fork & Clone

```bash
# Fork le projet sur GitHub / Fork the project on GitHub
git clone https://github.com/YOUR_USERNAME/manga-translator.git
cd manga-translator
```

### 2. Créer une branche / Create a Branch

Utilisez un nom descriptif pour votre branche :

Use a descriptive name for your branch:

```bash
git checkout -b feature/ma-nouvelle-fonctionnalité
# or
git checkout -b fix/correction-du-bug-xyz
```

### 3. Faire vos modifications / Make Your Changes

- Écrivez du code propre et commenté / Write clean and commented code
- Ajoutez des tests / Add tests
- Mettez à jour la documentation / Update documentation
- Vérifiez que les linters passent / Check that linters pass

### 4. Commit

Utilisez des messages de commit clairs et descriptifs :

Use clear and descriptive commit messages:

```bash
git add .
git commit -m "feat(ai-worker): ajout de la détection de bulles rondes"
# or
git commit -m "fix(backend): correction du timeout sur les requêtes longues"
```

Convention recommandée / Recommended convention:

- `feat`: Nouvelle fonctionnalité / New feature
- `fix`: Correction de bug / Bug fix
- `docs`: Documentation seulement / Documentation only
- `style`: Formatage, point-virgules, etc. / Formatting, semicolons, etc.
- `refactor`: Refactorisation du code / Code refactoring
- `test`: Ajout de tests / Adding tests
- `chore`: Maintenance, dépendances / Maintenance, dependencies

### 5. Push & Pull Request

```bash
git push origin feature/ma-nouvelle-fonctionnalité
```

Puis ouvrez une Pull Request sur GitHub avec :

Then open a Pull Request on GitHub with:

- **Titre clair** / **Clear title**
- **Description détaillée** de ce que vous avez fait et pourquoi / **Detailed description** of what you did and why
- **Référence à l'Issue** si applicable (ex: `Closes #42`) / **Reference to the Issue** if applicable (e.g., `Closes #42`)
- **Screenshots** ou exemples si c'est une modification visuelle / **Screenshots** or examples if it's a visual change

### 6. Review

- **P4ST4S** (Antoine Rospars) reviewera personnellement votre PR / will personally review your PR
- Vous pouvez aussi demander un review à **GitHub Copilot** si vous le souhaitez / You can also request a review from **GitHub Copilot** if you want
- Soyez ouvert aux feedbacks et prêt à faire des modifications / Be open to feedback and ready to make changes
- La review peut prendre quelques jours, soyez patient / The review may take a few days, be patient

---

## 🏗️ Setup de développement / Development Setup

### Prérequis / Prerequisites

- **Python 3.11+** avec CUDA (pour l'AI Worker) / with CUDA (for AI Worker)
- **Go 1.22+** (pour le Backend)
- **Node.js 18+** et **pnpm** (pour le Frontend)
- **Docker** (optionnel mais recommandé / optional but recommended)
- **Redis** (pour le Backend)
- **PostgreSQL** (pour le Backend)

### Installation locale / Local Installation

Consultez le [README.md](README.md) principal pour les instructions d'installation complètes.

Refer to the main [README.md](README.md) for complete installation instructions.

---

## 🤝 Code de Conduite / Code of Conduct

Nous n'avons pas de code de conduite formel, mais nous attendons de tous les contributeurs qu'ils :

We don't have a formal code of conduct, but we expect all contributors to:

- 🙏 **Restent courtois et respectueux** / **Stay courteous and respectful**
- 💡 **Respectent le niveau de compétence de chacun** / **Respect everyone's skill level**
- 🌈 **Soient ouverts aux différentes approches** / **Be open to different approaches**
- 🎯 **Se concentrent sur le code et les idées, pas sur les personnes** / **Focus on the code and ideas, not on people**

Les comportements irrespectueux, harcelants ou discriminatoires ne seront pas tolérés.

Disrespectful, harassing, or discriminatory behavior will not be tolerated.

---

## ❓ Questions ?

Si vous avez des questions ou besoin d'aide :

If you have questions or need help:

- 📝 Ouvrez une Issue avec le tag `question`
- 💬 Commentez sur une PR ou Issue existante

---

## 🎉 Remerciements / Acknowledgments

Merci à tous les contributeurs qui aident à améliorer ce projet !

Thanks to all contributors who help improve this project!

---

**Happy Coding! / Bon Code!** 🚀
