# 📚 Cahier des Charges - API Bibliothèque Personnelle Multi-utilisateurs

## 🎯 Objectif du Projet
Développer une API RESTful en Go permettant à des utilisateurs de gérer leurs bibliothèques personnelles de livres avec authentification et isolation des données.

## 👥 Acteurs du Système
- **Utilisateur non authentifié** : Peut s'inscrire et se connecter
- **Utilisateur authentifié** : Peut gérer sa bibliothèque personnelle et sa liste de souhaits

## 📋 Fonctionnalités Requises

### 1. **Gestion d'Utilisateurs**
- [ ] **Inscription** : Création de compte avec email et mot de passe
- [ ] **Connexion** : Authentification avec génération de token JWT
- [ ] **Profil utilisateur** : Consultation et modification du profil
- [ ] **Changement de mot de passe**
- [ ] **Déconnexion** : Invalidation du token

### 2. **Gestion des Livres (Bibliothèque Principale)**
- [ ] **Ajouter un livre** à sa bibliothèque
- [ ] **Modifier les informations** d'un livre existant
- [ ] **Supprimer un livre** de sa bibliothèque
- [ ] **Marquer comme lu/non lu**
- [ ] **Ajouter une note** (1-5 étoiles) et un commentaire personnel
- [ ] **Définir des dates** (début et fin de lecture)
- [ ] **Suivi de lecture** (en cours, à lire, terminé)

### 3. **Gestion de la Liste de Souhaits**
- [ ] **Ajouter un livre** à la liste de souhaits
- [ ] **Retirer un livre** de la liste de souhaits
- [ ] **Déplacer un livre** de la liste de souhaits vers la bibliothèque principale
- [ ] **Prioriser les souhaits** avec un niveau de priorité (faible, moyen, élevé)
- [ ] **Ajouter des notes** sur pourquoi ce livre est souhaité
- [ ] **Marquer comme acquis** lors de l'ajout à la bibliothèque

### 4. **Recherche et Filtrage**
- [ ] **Lister tous ses livres** avec pagination
- [ ] **Lister sa liste de souhaits** avec pagination
- [ ] **Recherche texte** par titre et auteur (dans bibliothèque ET souhaits)
- [ ] **Filtrer par statut de lecture** (à lire, en cours, terminé)
- [ ] **Filtrer par note** (1 à 5 étoiles)
- [ ] **Filtrer par priorité** dans la liste de souhaits
- [ ] **Trier par** : date d'ajout, titre, auteur, note, priorité
- [ ] **Recherche avancée** avec combinaison de filtres

### 5. **Fonctionnalités Avancées**
- [ ] **Statistiques personnelles** : 
  - Nombre total de livres
  - Nombre de livres lus/en cours/à lire
  - Nombre de livres dans la liste de souhaits
  - Moyenne des notes
  - Auteur le plus lu
  - Priorité moyenne des souhaits
- [ ] **Export de données** : Export de sa bibliothèque et liste de souhaits au format JSON
- [ ] **Import de livres** via ISBN (intégration avec API externe)
- [ ] **Suggestions** : Livres populaires parmi les souhaits des autres utilisateurs (anonymisées)

## 🔐 Sécurité et Authentification
- **JWT tokens** avec expiration courte (15-30 minutes)
- **Refresh tokens** pour maintenir la session
- **Hashage des mots de passe** avec bcrypt
- **Middleware d'authentification** sur toutes les routes protégées
- **Validation stricte** des données d'entrée
- **Isolation des données** : un utilisateur ne peut accéder qu'à ses propres livres et souhaits

## 🌐 API Endpoints Détaillés

### Authentification
```
POST /api/auth/register     # Création de compte
POST /api/auth/login        # Connexion
POST /api/auth/refresh      # Renouvellement du token
POST /api/auth/logout       # Déconnexion
GET  /api/auth/me           # Profil utilisateur
PUT  /api/auth/password     # Changement mot de passe
```

### Gestion des Livres (Bibliothèque)
```
GET    /api/books           # Liste avec filtres/pagination
GET    /api/books/{id}      # Détails d'un livre
POST   /api/books           # Ajouter un livre
PUT    /api/books/{id}      # Modifier un livre
DELETE /api/books/{id}      # Supprimer un livre
GET    /api/books/stats     # Statistiques personnelles
POST   /api/books/import    # Importer par ISBN
GET    /api/books/export    # Exporter sa bibliothèque
```

### Gestion de la Liste de Souhaits
```
GET    /api/wishlist        # Liste des souhaits avec filtres
GET    /api/wishlist/{id}   # Détails d'un souhait
POST   /api/wishlist        # Ajouter un livre à la liste de souhaits
PUT    /api/wishlist/{id}   # Modifier un souhait (priorité, notes)
DELETE /api/wishlist/{id}   # Retirer un livre de la liste de souhaits
POST   /api/wishlist/{id}/move-to-books  # Déplacer vers la bibliothèque
GET    /api/wishlist/stats  # Statistiques de la liste de souhaits
```

## 💾 Base de Données
- **PostgreSQL** comme base de données principale
- **Tables principales** : users, books, wishlist_items
- **Relations** : utilisateurs → livres (one-to-many), utilisateurs → souhaits (one-to-many)
- **Champs wishlist** : priorité, notes de souhait, date d'ajout, date d'acquisition
- **Index** sur les champs de recherche fréquents

## 🚀 Contraintes Techniques
- **Langage** : Go (avec framework web comme Gin ou Echo)
- **Base de données** : PostgreSQL
- **Authentification** : JWT
- **API** : RESTful JSON
- **Séparation** stricte des données par utilisateur
- **Gestion d'erreurs** appropriée avec codes HTTP significatifs

## 📦 Livrables
- Code source de l'API Go
- Scripts de configuration de la base PostgreSQL
- Documentation des endpoints API
- Instructions d'installation et déploiement

La liste de souhaits ajoute une dimension intéressante au projet tout en restant gérable pour l'apprentissage de Go et PostgreSQL.
