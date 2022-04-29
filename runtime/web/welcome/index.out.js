(() => {
  // ../../../node_modules/js-confetti/dist/es/index.js
  function _classCallCheck(instance, Constructor) {
    if (!(instance instanceof Constructor)) {
      throw new TypeError("Cannot call a class as a function");
    }
  }
  function _defineProperties(target, props) {
    for (var i = 0; i < props.length; i++) {
      var descriptor = props[i];
      descriptor.enumerable = descriptor.enumerable || false;
      descriptor.configurable = true;
      if ("value" in descriptor)
        descriptor.writable = true;
      Object.defineProperty(target, descriptor.key, descriptor);
    }
  }
  function _createClass(Constructor, protoProps, staticProps) {
    if (protoProps)
      _defineProperties(Constructor.prototype, protoProps);
    if (staticProps)
      _defineProperties(Constructor, staticProps);
    return Constructor;
  }
  function normalizeComputedStyleValue(string) {
    return +string.replace(/px/, "");
  }
  function fixDPR(canvas) {
    var dpr = window.devicePixelRatio;
    var computedStyles = getComputedStyle(canvas);
    var width = normalizeComputedStyleValue(computedStyles.getPropertyValue("width"));
    var height = normalizeComputedStyleValue(computedStyles.getPropertyValue("height"));
    canvas.setAttribute("width", (width * dpr).toString());
    canvas.setAttribute("height", (height * dpr).toString());
  }
  function generateRandomNumber(min, max) {
    var fractionDigits = arguments.length > 2 && arguments[2] !== void 0 ? arguments[2] : 0;
    var randomNumber = Math.random() * (max - min) + min;
    return Math.floor(randomNumber * Math.pow(10, fractionDigits)) / Math.pow(10, fractionDigits);
  }
  function generateRandomArrayElement(arr) {
    return arr[generateRandomNumber(0, arr.length)];
  }
  var FREE_FALLING_OBJECT_ACCELERATION = 125e-5;
  var MIN_DRAG_FORCE_COEFFICIENT = 5e-4;
  var MAX_DRAG_FORCE_COEFFICIENT = 9e-4;
  var ROTATION_SLOWDOWN_ACCELERATION = 1e-5;
  var INITIAL_SHAPE_RADIUS = 6;
  var INITIAL_EMOJI_SIZE = 80;
  var MIN_INITIAL_CONFETTI_SPEED = 0.9;
  var MAX_INITIAL_CONFETTI_SPEED = 1.7;
  var MIN_FINAL_X_CONFETTI_SPEED = 0.2;
  var MAX_FINAL_X_CONFETTI_SPEED = 0.6;
  var MIN_INITIAL_ROTATION_SPEED = 0.03;
  var MAX_INITIAL_ROTATION_SPEED = 0.07;
  var MIN_CONFETTI_ANGLE = 15;
  var MAX_CONFETTI_ANGLE = 82;
  var MAX_CONFETTI_POSITION_SHIFT = 150;
  var SHAPE_VISIBILITY_TRESHOLD = 100;
  var DEFAULT_CONFETTI_NUMBER = 250;
  var DEFAULT_EMOJIS_NUMBER = 40;
  var DEFAULT_CONFETTI_COLORS = ["#fcf403", "#62fc03", "#f4fc03", "#03e7fc", "#03fca5", "#a503fc", "#fc03ad", "#fc03c2"];
  function getWindowWidthCoefficient(canvasWidth) {
    var HD_SCREEN_WIDTH = 1920;
    return Math.log(canvasWidth) / Math.log(HD_SCREEN_WIDTH);
  }
  var ConfettiShape = /* @__PURE__ */ function() {
    function ConfettiShape2(args) {
      _classCallCheck(this, ConfettiShape2);
      var initialPosition = args.initialPosition, direction = args.direction, confettiRadius = args.confettiRadius, confettiColors = args.confettiColors, emojis = args.emojis, emojiSize = args.emojiSize, canvasWidth = args.canvasWidth;
      var randomConfettiSpeed = generateRandomNumber(MIN_INITIAL_CONFETTI_SPEED, MAX_INITIAL_CONFETTI_SPEED, 3);
      var initialSpeed = randomConfettiSpeed * getWindowWidthCoefficient(canvasWidth);
      this.confettiSpeed = {
        x: initialSpeed,
        y: initialSpeed
      };
      this.finalConfettiSpeedX = generateRandomNumber(MIN_FINAL_X_CONFETTI_SPEED, MAX_FINAL_X_CONFETTI_SPEED, 3);
      this.rotationSpeed = emojis.length ? 0.01 : generateRandomNumber(MIN_INITIAL_ROTATION_SPEED, MAX_INITIAL_ROTATION_SPEED, 3) * getWindowWidthCoefficient(canvasWidth);
      this.dragForceCoefficient = generateRandomNumber(MIN_DRAG_FORCE_COEFFICIENT, MAX_DRAG_FORCE_COEFFICIENT, 6);
      this.radius = {
        x: confettiRadius,
        y: confettiRadius
      };
      this.initialRadius = confettiRadius;
      this.rotationAngle = direction === "left" ? generateRandomNumber(0, 0.2, 3) : generateRandomNumber(-0.2, 0, 3);
      this.emojiSize = emojiSize;
      this.emojiRotationAngle = generateRandomNumber(0, 2 * Math.PI);
      this.radiusYUpdateDirection = "down";
      var angle = direction === "left" ? generateRandomNumber(MAX_CONFETTI_ANGLE, MIN_CONFETTI_ANGLE) * Math.PI / 180 : generateRandomNumber(-MIN_CONFETTI_ANGLE, -MAX_CONFETTI_ANGLE) * Math.PI / 180;
      this.absCos = Math.abs(Math.cos(angle));
      this.absSin = Math.abs(Math.sin(angle));
      var positionShift = generateRandomNumber(-MAX_CONFETTI_POSITION_SHIFT, 0);
      var shiftedInitialPosition = {
        x: initialPosition.x + (direction === "left" ? -positionShift : positionShift) * this.absCos,
        y: initialPosition.y - positionShift * this.absSin
      };
      this.currentPosition = Object.assign({}, shiftedInitialPosition);
      this.initialPosition = Object.assign({}, shiftedInitialPosition);
      this.color = emojis.length ? null : generateRandomArrayElement(confettiColors);
      this.emoji = emojis.length ? generateRandomArrayElement(emojis) : null;
      this.createdAt = new Date().getTime();
      this.direction = direction;
    }
    _createClass(ConfettiShape2, [{
      key: "draw",
      value: function draw(canvasContext) {
        var currentPosition = this.currentPosition, radius = this.radius, color = this.color, emoji = this.emoji, rotationAngle = this.rotationAngle, emojiRotationAngle = this.emojiRotationAngle, emojiSize = this.emojiSize;
        var dpr = window.devicePixelRatio;
        if (color) {
          canvasContext.fillStyle = color;
          canvasContext.beginPath();
          canvasContext.ellipse(currentPosition.x * dpr, currentPosition.y * dpr, radius.x * dpr, radius.y * dpr, rotationAngle, 0, 2 * Math.PI);
          canvasContext.fill();
        } else if (emoji) {
          canvasContext.font = "".concat(emojiSize, "px serif");
          canvasContext.save();
          canvasContext.translate(dpr * currentPosition.x, dpr * currentPosition.y);
          canvasContext.rotate(emojiRotationAngle);
          canvasContext.textAlign = "center";
          canvasContext.fillText(emoji, 0, 0);
          canvasContext.restore();
        }
      }
    }, {
      key: "updatePosition",
      value: function updatePosition(iterationTimeDelta, currentTime) {
        var confettiSpeed = this.confettiSpeed, dragForceCoefficient = this.dragForceCoefficient, finalConfettiSpeedX = this.finalConfettiSpeedX, radiusYUpdateDirection = this.radiusYUpdateDirection, rotationSpeed = this.rotationSpeed, createdAt = this.createdAt, direction = this.direction;
        var timeDeltaSinceCreation = currentTime - createdAt;
        if (confettiSpeed.x > finalConfettiSpeedX)
          this.confettiSpeed.x -= dragForceCoefficient * iterationTimeDelta;
        this.currentPosition.x += confettiSpeed.x * (direction === "left" ? -this.absCos : this.absCos) * iterationTimeDelta;
        this.currentPosition.y = this.initialPosition.y - confettiSpeed.y * this.absSin * timeDeltaSinceCreation + FREE_FALLING_OBJECT_ACCELERATION * Math.pow(timeDeltaSinceCreation, 2) / 2;
        this.rotationSpeed -= this.emoji ? 1e-4 : ROTATION_SLOWDOWN_ACCELERATION * iterationTimeDelta;
        if (this.rotationSpeed < 0)
          this.rotationSpeed = 0;
        if (this.emoji) {
          this.emojiRotationAngle += this.rotationSpeed * iterationTimeDelta % (2 * Math.PI);
          return;
        }
        if (radiusYUpdateDirection === "down") {
          this.radius.y -= iterationTimeDelta * rotationSpeed;
          if (this.radius.y <= 0) {
            this.radius.y = 0;
            this.radiusYUpdateDirection = "up";
          }
        } else {
          this.radius.y += iterationTimeDelta * rotationSpeed;
          if (this.radius.y >= this.initialRadius) {
            this.radius.y = this.initialRadius;
            this.radiusYUpdateDirection = "down";
          }
        }
      }
    }, {
      key: "getIsVisibleOnCanvas",
      value: function getIsVisibleOnCanvas(canvasHeight) {
        return this.currentPosition.y < canvasHeight + SHAPE_VISIBILITY_TRESHOLD;
      }
    }]);
    return ConfettiShape2;
  }();
  function createCanvas() {
    var canvas = document.createElement("canvas");
    canvas.style.position = "fixed";
    canvas.style.width = "100%";
    canvas.style.height = "100%";
    canvas.style.top = "0";
    canvas.style.left = "0";
    canvas.style.zIndex = "1000";
    canvas.style.pointerEvents = "none";
    document.body.appendChild(canvas);
    return canvas;
  }
  function normalizeConfettiConfig(confettiConfig) {
    var _confettiConfig$confe = confettiConfig.confettiRadius, confettiRadius = _confettiConfig$confe === void 0 ? INITIAL_SHAPE_RADIUS : _confettiConfig$confe, _confettiConfig$confe2 = confettiConfig.confettiNumber, confettiNumber = _confettiConfig$confe2 === void 0 ? confettiConfig.confettiesNumber || (confettiConfig.emojis ? DEFAULT_EMOJIS_NUMBER : DEFAULT_CONFETTI_NUMBER) : _confettiConfig$confe2, _confettiConfig$confe3 = confettiConfig.confettiColors, confettiColors = _confettiConfig$confe3 === void 0 ? DEFAULT_CONFETTI_COLORS : _confettiConfig$confe3, _confettiConfig$emoji = confettiConfig.emojis, emojis = _confettiConfig$emoji === void 0 ? confettiConfig.emojies || [] : _confettiConfig$emoji, _confettiConfig$emoji2 = confettiConfig.emojiSize, emojiSize = _confettiConfig$emoji2 === void 0 ? INITIAL_EMOJI_SIZE : _confettiConfig$emoji2;
    if (confettiConfig.emojies)
      console.error("emojies argument is deprecated, please use emojis instead");
    if (confettiConfig.confettiesNumber)
      console.error("confettiesNumber argument is deprecated, please use confettiNumber instead");
    return {
      confettiRadius,
      confettiNumber,
      confettiColors,
      emojis,
      emojiSize
    };
  }
  var ConfettiBatch = /* @__PURE__ */ function() {
    function ConfettiBatch2(canvasContext) {
      var _this = this;
      _classCallCheck(this, ConfettiBatch2);
      this.canvasContext = canvasContext;
      this.shapes = [];
      this.promise = new Promise(function(completionCallback) {
        return _this.resolvePromise = completionCallback;
      });
    }
    _createClass(ConfettiBatch2, [{
      key: "getBatchCompletePromise",
      value: function getBatchCompletePromise() {
        return this.promise;
      }
    }, {
      key: "addShapes",
      value: function addShapes() {
        var _this$shapes;
        (_this$shapes = this.shapes).push.apply(_this$shapes, arguments);
      }
    }, {
      key: "complete",
      value: function complete() {
        var _a;
        if (this.shapes.length) {
          return false;
        }
        (_a = this.resolvePromise) === null || _a === void 0 ? void 0 : _a.call(this);
        return true;
      }
    }, {
      key: "processShapes",
      value: function processShapes(time, canvasHeight, cleanupInvisibleShapes) {
        var _this2 = this;
        var timeDelta = time.timeDelta, currentTime = time.currentTime;
        this.shapes = this.shapes.filter(function(shape) {
          shape.updatePosition(timeDelta, currentTime);
          shape.draw(_this2.canvasContext);
          if (!cleanupInvisibleShapes) {
            return true;
          }
          return shape.getIsVisibleOnCanvas(canvasHeight);
        });
      }
    }]);
    return ConfettiBatch2;
  }();
  var JSConfetti = /* @__PURE__ */ function() {
    function JSConfetti2() {
      var jsConfettiConfig = arguments.length > 0 && arguments[0] !== void 0 ? arguments[0] : {};
      _classCallCheck(this, JSConfetti2);
      this.activeConfettiBatches = [];
      this.canvas = jsConfettiConfig.canvas || createCanvas();
      this.canvasContext = this.canvas.getContext("2d");
      this.requestAnimationFrameRequested = false;
      this.lastUpdated = new Date().getTime();
      this.iterationIndex = 0;
      this.loop = this.loop.bind(this);
      requestAnimationFrame(this.loop);
    }
    _createClass(JSConfetti2, [{
      key: "loop",
      value: function loop() {
        this.requestAnimationFrameRequested = false;
        fixDPR(this.canvas);
        var currentTime = new Date().getTime();
        var timeDelta = currentTime - this.lastUpdated;
        var canvasHeight = this.canvas.offsetHeight;
        var cleanupInvisibleShapes = this.iterationIndex % 10 === 0;
        this.activeConfettiBatches = this.activeConfettiBatches.filter(function(batch) {
          batch.processShapes({
            timeDelta,
            currentTime
          }, canvasHeight, cleanupInvisibleShapes);
          if (!cleanupInvisibleShapes) {
            return true;
          }
          return !batch.complete();
        });
        this.iterationIndex++;
        this.queueAnimationFrameIfNeeded(currentTime);
      }
    }, {
      key: "queueAnimationFrameIfNeeded",
      value: function queueAnimationFrameIfNeeded(currentTime) {
        if (this.requestAnimationFrameRequested) {
          return;
        }
        if (this.activeConfettiBatches.length < 1) {
          return;
        }
        this.requestAnimationFrameRequested = true;
        this.lastUpdated = currentTime || new Date().getTime();
        requestAnimationFrame(this.loop);
      }
    }, {
      key: "addConfetti",
      value: function addConfetti() {
        var confettiConfig = arguments.length > 0 && arguments[0] !== void 0 ? arguments[0] : {};
        var _normalizeConfettiCon = normalizeConfettiConfig(confettiConfig), confettiRadius = _normalizeConfettiCon.confettiRadius, confettiNumber = _normalizeConfettiCon.confettiNumber, confettiColors = _normalizeConfettiCon.confettiColors, emojis = _normalizeConfettiCon.emojis, emojiSize = _normalizeConfettiCon.emojiSize;
        var canvasRect = this.canvas.getBoundingClientRect();
        var canvasWidth = canvasRect.width;
        var canvasHeight = canvasRect.height;
        var yPosition = canvasHeight * 5 / 7;
        var leftConfettiPosition = {
          x: 0,
          y: yPosition
        };
        var rightConfettiPosition = {
          x: canvasWidth,
          y: yPosition
        };
        var confettiGroup = new ConfettiBatch(this.canvasContext);
        for (var i = 0; i < confettiNumber / 2; i++) {
          var confettiOnTheRight = new ConfettiShape({
            initialPosition: leftConfettiPosition,
            direction: "right",
            confettiRadius,
            confettiColors,
            confettiNumber,
            emojis,
            emojiSize,
            canvasWidth
          });
          var confettiOnTheLeft = new ConfettiShape({
            initialPosition: rightConfettiPosition,
            direction: "left",
            confettiRadius,
            confettiColors,
            confettiNumber,
            emojis,
            emojiSize,
            canvasWidth
          });
          confettiGroup.addShapes(confettiOnTheRight, confettiOnTheLeft);
        }
        this.activeConfettiBatches.push(confettiGroup);
        this.queueAnimationFrameIfNeeded();
        return confettiGroup.getBatchCompletePromise();
      }
    }]);
    return JSConfetti2;
  }();
  var es_default = JSConfetti;

  // index.js
  var jsConfetti = new es_default();
  jsConfetti.addConfetti();
  var sse = new EventSource("http://0.0.0.0:35729");
  sse.addEventListener("message", () => {
    location.reload();
  });
})();
