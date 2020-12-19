======
tchart
======

Command line tool to visualize and inspect table based data

.. image:: docs/overview_tchart.jpg

Installation
============

.. code::

    go get github.com/s-westphal/go-tchart/cmd/tchart

Usage
=====

.. code::

    tchart --help
   
Input data is read from stdin.

Panels
------

Every panel can show data from several columns of the input data.  

On the left side of each panel "statsWidgets" show a summary of each column:

.. image:: docs/data_stats.jpg
    :width: 12%

The main chart on the right shows the joined data from all it's columns. The kind 
of chart being displayed is specified by the panel type.

**Panel/Chart Types**:

- Panel type ``[L]`` -> Line chart  (default)

    .. image:: docs/line_chart.jpg

- Panel type ``[S]`` -> Scatter plot 

    .. image:: docs/scatter_plot.jpg
    

- Panel type ``[P]`` -> Pie chart 

    .. image:: docs/pie_chart.jpg


The ``-p <panel-types>`` cli-option can be used to define the displayed panels.  
``<panel-types>`` is a string containing one character for every column of data,
where the n-th character corresponds to the n-th column.

Possible characters are:

- ``L``, ``S``, ``P``: create a new panel with the according type and display this colum's data.
- ``.``: add this column's data to the previous panel.
- ``x``: skip this column.



**Example:**

    Input Data:

    +--------+----------+-------+-----------+------------+
    | value  | f(value) |skip   | progress_a|  progress_b|
    +========+==========+=======+===========+============+
    |1       | 5        | foo   | 5         |7           |
    +--------+----------+-------+-----------+------------+
    |\...    |          |       |           |            |
    +--------+----------+-------+-----------+------------+


    .. code::

        cat input_data | tchart -p S.xL.


    -> ``S.xL.`` is mapped to the columns:

    .. code::

        "value"         ->  S  =>  create ScatterPlot-Panel
        "f(value)"      ->  .  =>  add to ScatterPlot-Panel
        "skip"          ->  x  =>  skip this column  
        "progress_a"    ->  L  =>  create LineChart-Panel
        "progress_b"    ->  .  =>  add to LineChart-Panel


    => Two panels are displayed, the first one showing a scatter plot with x-values 
    ``value`` and y-values ``f(value)``, the column ``random`` is skipped. The second 
    panel shows a line chart with data from columns ``progress_a`` and ``progress_b``.

**Notes:**

- By default every column is displayed in a separate line chart.
- Scatter plots must contain exactly 2 columns.



Examples
--------

- Scatter plot with 50 points displayed at maximum:

    .. code::

        seq 500 | awk 'BEGIN{OFS="\t"; print "rand","2*rand"}{x=$1/5; print rand(),2*rand}' | tchart -p S. -n 50


- Skip columns using ``x`` in panels option and load data fast:

    .. code::

        seq 500 | awk 'BEGIN{OFS="\t"; print "skip","sin(x)","skip","cos(x)"}{x=$1/5; print "foo",sin(x),"bar",cos(x)}' | tchart -p xLx. -s fast

- Use first column as labels with option ``-l first``:

    .. code::

        seq 500 | awk 'BEGIN{OFS="\t"; print "x","sin(x)","cos(x)"}{x=$1/5; print x,sin(x),cos(x)}' | tchart -p L. -l first

- Show different panels at once:

    .. code::

        seq 1500 | awk 'BEGIN{OFS="\t"; print "x","2*sin(2*x)","cos(x)","sin(x)","2*cos(x)","rand","3*rand","2*rand"}{x=$1/5; print x,2*sin(2*x),cos(x),sin(x),2*cos(x),rand(),3*rand(),2*rand()}' | tchart -p L.S.P.. -l first



License
=======

`MIT <http://opensource.org/licenses/MIT>`_
