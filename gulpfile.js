var gulp = require('gulp'),
    sass = require('gulp-sass'),
    bro = require('gulp-bro'),
    postcss = require('gulp-postcss');
    
gulp.task('build:css', function () {
  return (
    gulp.src('./src/css/**/*.scss')
    .pipe(sass.sync().on('error', sass.logError))
    .pipe(postcss())
    .pipe(gulp.dest('./dist/css'))
  )
});

gulp.task('build:image', function () {
  return (
    gulp.src('./src/img/**/*.png')
    .pipe(gulp.dest('./dist/img'))
  )
});

gulp.task('build:commonjs', function () {
  // take every commonJS module, browserify them and put them into ./dist
  return (
    gulp.src('./src/cjs-modules/*.js')
        .pipe(bro())
        .pipe(gulp.dest('./dist/js'))
  )
})

gulp.task('build:js', gulp.series('build:commonjs', function(){
  return (
    // take every JS script, and put them into ./dist
    gulp.src('./src/js/**/*.js')
        .pipe(gulp.dest('./dist/js'))
  )
}))

gulp.task('build', gulp.series('build:css', 'build:js', 'build:image'))