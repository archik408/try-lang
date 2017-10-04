from __future__ import unicode_literals

from django.db import models

class Person(models.Model):
    """A model of a person."""
    id = models.CharField("Identifier", max_length=20, unique=True, primary_key=True)
    name = models.CharField("User's first name", max_length=20)
    data = models.CharField("User's data", max_length=20)