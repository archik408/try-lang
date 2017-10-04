from django.core import serializers
from django.http import HttpResponse
from django.db.models import Q

from main_app import models

def list(request):
    """Get all persons."""
    if request.method == 'GET':
        persons = models.Person.objects.all()
        return _reply(persons)
    else:
        return _reply(status=400)

def createOrModify(request, id):
    """Get or create/update/patch specific person"""
    if request.method == 'GET':
        try:
            person = models.Person.objects.get(id=id)
            return _reply([person])
        except:
            return _reply(status=404)
    elif request.method == 'POST' or request.method == 'PUT' or request.method == 'PATCH':
        for person in serializers.deserialize("json", request.body):
            person.save()
        return _reply()
    else:
        return _reply(status=400)

def find(request):
    """Find specific person"""
    if request.method == 'GET':
        text = request.GET.get('text', '')
        persons = models.Person.objects.filter(Q(name=text) | Q(data=text))
        return _reply(persons)
    else:
        return _reply(status=400)

def _reply(data = [], status=200):
    return HttpResponse(serializers.serialize("json", data), content_type="application/json", status=status)