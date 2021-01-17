$(function(){
  // alertTimeout it's the default timeout
  // to show the alerts, 5s
  var alertTimeout = 5000
  var currentTimout, currentAlert;

  // catchError it's a global function that will call showAlert
  // with the response error
  var catchError = function(collection, response, options) {
    showAlert(response.responseJSON.error, "danger")
  }

  var Notification = Backbone.Model.extend({
    attributes: {
      "id": "",
      "title": "",
      "url": "",
      "scope": "",
      "updated_at": "",
    },
  });

  var NotificationList = Backbone.Collection.extend({
    model: Notification,

    comparator: 'updated_at',

    parse: function(response, options) {
      return response.data;
    },

    url: function() {
      return this.source.url() + "/notifications"
    },

    groupByScopes: function() {
      var gps = this.groupBy(function(e) {
        return e.get("scope");
      });

      var res = [];

      _.each(gps, function(notifications, scope) {
        var ne =_.sortBy(notifications, function(event) {
          return event.get("updated_at");
        });
        ne = ne.reverse();
        res.push({
          scope: scope,
          notifications: ne,
          updated_at: ne[0].get("updated_at"),
        });
      });

      return _.sortBy(res, function(item) {
        return item.updated_at;
      }).reverse();
    },
  });

  var Source = Backbone.Model.extend({
    idAttribute: "canonical",
    attributes: {
      "canonical": "",
      "name": "",
      "active": false,
    },
    urlRoot: "/api/sources",
    initialize: function(options) {
      this.notifications = new NotificationList();
      this.notifications.source = this;
    },
  });

  var Alert = Backbone.Model.extend({
    attributes: {
      "text": "",
      "type": "",
    }
  })

  var SourceList = Backbone.Collection.extend({
    model: Source,
    url: '/api/sources',

    initialize: function() {
      this.active = undefined;
    },

    parse: function(response, options) {
      return response.data;
    },

    toggleActive: function(current) {
      if (this.active != undefined) this.active.set('active', false);

      this.active = this.get(current);
      this.active.set('active', true);
    },
  });

  var NotificationsView = Backbone.View.extend({
    template: _.template($('#notifications-view-tmpl').html()),
    events: {
      "click .scope-title": "toggleNotifications",
      "click button#refresh": "refreshNotifications",
    },
    initialize: function() {
      this.listenTo(this.collection, "reset", this.render);
      showAlert("Refreshing ...", "warning");
      this.collection.fetch({
        reset: true,
        success: function(){showAlert("Refreshed!", "success")},
        error: catchError,
      });
    },
    render: function() {
      this.$el.html(this.template({ notifications: this.collection }));
      return this;
    },
    toggleNotifications: function(e) {
      this.$el.find("#notifications-"+e.currentTarget.id).toggle("fast")
    },
    refreshNotifications: function(e) {
      showAlert("Refreshing ...", "warning");
      this.collection.fetch({
        reset: true,
        success: function(){showAlert("Refreshed!", "success")},
        error: catchError,
      })
    },
  })

  var SourceView = Backbone.View.extend({
    tagName: 'a',
    template: _.template($('#source-view-tmpl').html()),
    className: 'list-group-item list-group-item-action',
    events: {
      "click": "goToSource",
    },
    initialize: function() {
      this.listenTo(this.model, "change", this.render);
    },
    attributes: function(){
      return {
        'href': this.model.get('canonical'),
      };
    },
    goToSource: function(e) {
      e.preventDefault();
      router.navigate("/"+this.model.get("canonical"), {trigger: true})
    },
    render: function() {
      if (this.model.get('active')) {
        $(this.el).addClass('active').html(this.template(this.model.toJSON()));
      } else {
        $(this.el).removeClass('active').html(this.template(this.model.toJSON()));
      }
      return this;
    },
  });

  var AlertView = Backbone.View.extend({
    attributes: function() {
      return {
        class: 'alert alert-' + this.model.get("type"),
      };
    },
    template: _.template($('#alert-view-tmpl').html()),
    model: Alert,
    render: function() {
      this.$el.html(this.template(this.model.toJSON()));
      return this;
    }
  })

  // showAlert prints an alert with the
  // given text
  var showAlert = function(text, type) {
    if (currentAlert) currentAlert.remove();
    var a = new AlertView({model: new Alert({text: text, type: type})});
    currentAlert = a;
    $("#alert").html(a.render().el);
    if (currentTimout) {
      clearTimeout(currentTimout)
    }
    currentTimout = setTimeout(function() {
      currentAlert.remove();
    }, alertTimeout);
  }

  var Sources = new SourceList();

  var AppView = Backbone.View.extend({
    initialize: function() {
      this.$content = this.$("#content");
      this.$sidebar = this.$("#sidebar");
      this.sidebars = []

      this.listenTo(Sources, "sync", this.renderSidebar)
    },

    setContent: function(view) {
      var content = this.content;
      if (content) content.remove();
      content = this.content = view;
      this.$content.html(content.render().el);
    },

    renderSidebar: function() {
      var that = this
      var sidebars = that.sidebars
      if (sidebars.length != 0) _.each(sidebars, function(sv) { sv.remove() })
      that.sidebars = [];
      Sources.each(function(source) {
        var sv = new SourceView({model: source});
        that.sidebars.push(sv)
        that.$sidebar.append(sv.render().el);
      });
    },
  });


  var Router = Backbone.Router.extend({
    routes: {
      ":sourceCan": "sourceRender",
      "*all":       "renderHome",
    },
    initialize: function() {
      this.layout = new AppView({
        el: 'body',
      });
    },

    sourceRender: function(sourceCan) {
      var s = Sources.get(sourceCan);
      Sources.toggleActive(sourceCan);

      var nv = new NotificationsView({collection: s.notifications});
      this.layout.setContent(nv)
    },

    renderHome: function() {
      var src = Sources.first()
      router.navigate("/"+src.get("canonical"), {trigger: true})
    },
  });


  var router;

  // As we need the Sources to be fetched before doing anything we
  // initialize the router once they are fetched
  Sources.fetch({
    reset: true,
    success: function() {
      router = new Router();

      Backbone.ajax = function(request) {
        request = _({ contentType: 'application/json' }).defaults(request);
        return Backbone.$.ajax.call(Backbone.$, request);
      };

      Backbone.history.start({pushState: true});
    },
    error: catchError,
  });

})
