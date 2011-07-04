import json
import time

import datastore
import idmanager


class Paste(object):
  '''Abstraction for storing and retrieving pastes.'''

  class SaveStates:
    CLEAN = 0
    DIRTY = 1

  def __init__(self, ds, id=None):
    '''Constructor. If id is specified, the given paste id is retrieved.
    Otherwise, a new paste ID is allocated.

    Args:
      ds: datastore.Datastore object
      id (optional): string ID of paste to retrieve
    '''
    self._ds = ds

    # Clear all attributes.
    self._lexer_alias = None
    self._ttl = None
    self._ip_address = None
    self._content = None
    self._created = 0

    if id is None:
      mgr = idmanager.IDManager(ds)
      self._id = mgr.new_id()
      self._created = int(time.time())
      self._save_state = self.SaveStates.DIRTY
    else:
      self._id = id
      self._load()

  def _load(self):
    '''Loads attributes from the datastore.

    Raises:
      KeyError on nonexistent ID
    '''
    # Save the metadata.
    md_key = 'paste:%s' % self._id
    md_json = self._ds.value(md_key)
    if md_json is None:
      raise KeyError, 'no such ID'
    md = json.loads(md_json)

    self._lexer_alias = md.get('lexer', None)
    self._ttl = md.get('ttl', None)
    self._ip_address = md.get('ip_address', None)
    self._created = md.get('created', 0)
    self._content = md.get('content', None)
    self._save_state = self.SaveStates.CLEAN

  def _save(self):
    '''Saves attributes to the datastore.'''
    # Save the metadata.
    md_key = 'paste:%s' % self._id
    md = {'ttl': self._ttl,
          'lexer': self._lexer_alias,
          'ip_address': self._ip_address,
          'created': self._created,
          'content': self._content}
    self._ds.value_is(md_key, json.dumps(md))

  def save_state(self):
    '''Returns the save state (e.g., clean or dirty).

    Returns:
      Paste.SaveStates.* class constant
    '''
    return self._save_state

  def save_state_is(self, state):
    '''Changes the save state and does any necessary action (i.e., committing
    unsaved data) in order to accomplish this.

    Args:
      state: new save state

    Raises:
      ValueError: new state is invalid (e.g., clean -> dirty is not valid
      through this mechanism)
    '''
    if self._save_state == state:
      # no-op
      return

    if (state != self.SaveStates.CLEAN and
        state != self.SaveStates.DIRTY):
      raise ValueError, 'unknown state'

    if (self._save_state == self.SaveStates.CLEAN and
        state == self.SaveStates.DIRTY):
      # clean -> dirty
      raise ValueError, ('cannot change state from clean to dirty without '
                         'another operation')
    else:
      # dirty -> clean
      self._save()
      self._save_state = state

  def lexer_alias(self):
    '''Returns the lexer alias for this paste or None if not set.

    Returns:
      string (e.g., 'py') or None
    '''
    return self._lexer_alias

  def lexer_alias_is(self, alias):
    '''Sets the lexer alias for this paste.

    Args:
      alias: string lexer alias (e.g., 'py' or 'pl')
    '''
    self._lexer_alias = alias
    self._save_state = self.SaveStates.DIRTY

  def ttl(self):
    '''Returns the TTL for this paste or None if not set.

    Returns:
      ttl in seconds (integer)
    '''
    return self._ttl

  def ttl_is(self, ttl):
    '''Sets the ttl for this paste.

    Args:
      ttl: ttl integer
    '''
    self._ttl = ttl
    self._save_state = self.SaveStates.DIRTY

  def ip_address(self):
    '''Returns the IP address 'owner' associated with this paste.

    Returns:
      IP address string
    '''
    return self._ip_address

  def ip_address_is(self, ip_address):
    '''Sets the IP address 'owner'.

    Args:
      ip_address: IP address string
    '''
    self._ip_address = ip_address
    self._save_state = self.SaveStates.DIRTY

  def content(self):
    '''Returns the paste content.

    Returns:
      unicode string
    '''
    return self._content

  def content_is(self, content):
    '''Sets the paste content.

    Args:
      content: string content
    '''
    self._content = content
    self._save_state = self.SaveStates.DIRTY

  def id(self):
    '''Returns the ID for this paste.

    Returns:
      id string
    '''
    return self._id

  def created(self):
    '''Returns the epoch at which this paste was created.

    Returns:
      epoch value
    '''
    return self._created
